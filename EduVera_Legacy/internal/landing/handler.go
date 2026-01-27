package landing

import (
	"context"
	"eduvera/internal/platform/database"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"

	"golang.org/x/crypto/bcrypt"
)

// Home menangani request ke halaman utama landing page
func Home(w http.ResponseWriter, r *http.Request) {
	tmplPath := filepath.Join("web", "templates", "landing", "index.html")

	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		http.Error(w, "Gagal memuat template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Title": "EduVera - Solusi Sekolah & Pesantren Terpadu",
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Gagal render template: "+err.Error(), http.StatusInternalServerError)
	}
}

// Register menangani halaman pendaftaran tenant baru
func Register(w http.ResponseWriter, r *http.Request) {
	// Jika method POST, proses penyimpanan data
	if r.Method == http.MethodPost {
		fullname := r.FormValue("fullname")
		email := r.FormValue("email")
		wa := r.FormValue("wa")
		password := r.FormValue("password")

		// 1. Hash Password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Gagal enkripsi password", http.StatusInternalServerError)
			return
		}

		// 2. Simpan ke Database
		// Saat registrasi awal, tenant_id NULL dulu. Nanti diupdate setelah step Create Tenant.
		// Menggunakan QueryRow untuk mengambil ID user yang baru dibuat
		query := `INSERT INTO users (email, password_hash, full_name, phone, role)
				  VALUES ($1, $2, $3, $4, 'ADMIN_SEKOLAH') RETURNING id`

		var userID string
		err = database.DB.QueryRow(context.Background(), query, email, string(hashedPassword), fullname, wa).Scan(&userID)
		if err != nil {
			// Handle duplicate email error secara kasar dulu
			http.Error(w, "Gagal menyimpan data user (Email mungkin sudah terdaftar): "+err.Error(), http.StatusInternalServerError)
			return
		}

		// 3. Redirect ke Step 2 (Pilih Jenis Lembaga)
		// Kita bawa user_id di query param untuk context step selanjutnya
		http.Redirect(w, r, "/onboarding/step-2?user_id="+userID, http.StatusSeeOther)
		return
	}

	// Jika GET, tampilkan form
	tmplPath := filepath.Join("web", "templates", "landing", "register.html")

	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		http.Error(w, "Gagal memuat template register: "+err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Title": "Daftar EduVera - Mulai Kelola Sekolah Anda",
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Gagal render template register: "+err.Error(), http.StatusInternalServerError)
	}
}

// OnboardingStep2 menangani pemilihan jenis lembaga
func OnboardingStep2(w http.ResponseWriter, r *http.Request) {
	// Ambil user_id dari context
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Redirect(w, r, "/register", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodPost {
		moduleType := r.FormValue("type") // school / pesantren

		// Redirect ke Step 3 (Subdomain) dengan membawa data
		http.Redirect(w, r, "/onboarding/step-3?user_id="+userID+"&type="+moduleType, http.StatusSeeOther)
		return
	}

	tmplPath := filepath.Join("web", "templates", "landing", "step2.html")
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		http.Error(w, "Gagal memuat template step2: "+err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Title":  "Pilih Jenis Lembaga - EduVera",
		"UserID": userID,
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Gagal render template step2: "+err.Error(), http.StatusInternalServerError)
	}
}

// OnboardingStep3 menangani input subdomain dan pembuatan tenant
func OnboardingStep3(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	moduleType := r.URL.Query().Get("type")

	if userID == "" || moduleType == "" {
		http.Redirect(w, r, "/register", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodPost {
		schoolName := r.FormValue("school_name")
		subdomain := r.FormValue("subdomain")

		// Tentukan module flag berdasarkan pilihan Step 2
		hasSchool := false
		hasPesantren := false
		var level string

		if moduleType == "school" {
			hasSchool = true
			level = "SMA" // Default level, nanti bisa diubah di settings
		} else if moduleType == "pesantren" {
			hasPesantren = true
			level = "PESANTREN"
		}

		// 1. Insert Tenant
		// Menggunakan transaction agar atomik (Create Tenant + Link User)
		tx, err := database.DB.Begin(context.Background())
		if err != nil {
			http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer tx.Rollback(context.Background())

		var tenantID string
		queryTenant := `INSERT INTO tenants (name, subdomain, level, has_school_module, has_pesantren_module)
						VALUES ($1, $2, $3, $4, $5) RETURNING id`

		err = tx.QueryRow(context.Background(), queryTenant, schoolName, subdomain, level, hasSchool, hasPesantren).Scan(&tenantID)
		if err != nil {
			// Handle duplicate subdomain error
			http.Error(w, "Gagal membuat tenant (Subdomain mungkin sudah dipakai): "+err.Error(), http.StatusInternalServerError)
			return
		}

		// 2. Update User (Link to Tenant)
		queryUser := `UPDATE users SET tenant_id = $1 WHERE id = $2`
		_, err = tx.Exec(context.Background(), queryUser, tenantID, userID)
		if err != nil {
			http.Error(w, "Gagal update user: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Commit transaction
		err = tx.Commit(context.Background())
		if err != nil {
			http.Error(w, "Gagal commit transaksi: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// 3. Redirect ke Step 4 (Setup Keuangan)
		http.Redirect(w, r, "/onboarding/step-4?tenant_id="+tenantID, http.StatusSeeOther)
		return
	}

	tmplPath := filepath.Join("web", "templates", "landing", "step3.html")
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		http.Error(w, "Gagal memuat template step3: "+err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Title":  "Setup Subdomain - EduVera",
		"UserID": userID,
		"Type":   moduleType,
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Gagal render template step3: "+err.Error(), http.StatusInternalServerError)
	}
}

// OnboardingStep4 menangani setup data keuangan (Rekening Penarikan)
func OnboardingStep4(w http.ResponseWriter, r *http.Request) {
	tenantID := r.URL.Query().Get("tenant_id")
	if tenantID == "" {
		http.Redirect(w, r, "/register", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodPost {
		bankName := r.FormValue("bank_name")
		accNumber := r.FormValue("account_number")
		accHolder := r.FormValue("account_holder")

		// Update Tenant dengan data rekening
		query := `UPDATE tenants SET
			withdrawal_bank_name = $1,
			withdrawal_account_number = $2,
			withdrawal_account_holder = $3
			WHERE id = $4`

		_, err := database.DB.Exec(context.Background(), query, bankName, accNumber, accHolder, tenantID)
		if err != nil {
			http.Error(w, "Gagal update data keuangan: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Redirect ke Final Step (Summary / Payment)
		http.Redirect(w, r, "/onboarding/step-5?tenant_id="+tenantID, http.StatusSeeOther)
		return
	}

	tmplPath := filepath.Join("web", "templates", "landing", "step4.html")
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		http.Error(w, "Gagal memuat template step4: "+err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Title":    "Setup Keuangan - EduVera",
		"TenantID": tenantID,
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Gagal render template step4: "+err.Error(), http.StatusInternalServerError)
	}
}

// OnboardingStep5 menangani payment dan aktivasi
func OnboardingStep5(w http.ResponseWriter, r *http.Request) {
	tenantID := r.URL.Query().Get("tenant_id")
	if tenantID == "" {
		http.Redirect(w, r, "/register", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodPost {
		// SIMULASI PEMBAYARAN SUKSES
		// Update status tenant jadi ACTIVE (jika belum)
		// Sebenarnya default sudah ACTIVE, tapi ini untuk memastikan flow
		query := `UPDATE tenants SET status = 'ACTIVE' WHERE id = $1`
		_, err := database.DB.Exec(context.Background(), query, tenantID)
		if err != nil {
			http.Error(w, "Gagal aktivasi tenant: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Redirect ke Dashboard Utama (Selesai Onboarding)
		http.Redirect(w, r, "/dashboard?tenant_id="+tenantID, http.StatusSeeOther)
		return
	}

	// Ambil data tenant untuk ditampilkan di invoice ringkasan
	var name, subdomain string
	var hasSchool, hasPesantren bool

	query := `SELECT name, subdomain, has_school_module, has_pesantren_module FROM tenants WHERE id = $1`
	err := database.DB.QueryRow(context.Background(), query, tenantID).Scan(&name, &subdomain, &hasSchool, &hasPesantren)
	if err != nil {
		http.Error(w, "Gagal mengambil data tenant: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Hitung harga (Dummy Logic)
	var price int
	var packageName string

	if hasSchool && hasPesantren {
		price = 1200000
		packageName = "Hybrid Enterprise"
	} else if hasPesantren {
		price = 750000
		packageName = "Pesantren Pro"
	} else {
		price = 500000
		packageName = "Sekolah Dasar"
	}

	tmplPath := filepath.Join("web", "templates", "landing", "step5.html")
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		http.Error(w, "Gagal memuat template step5: "+err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Title":          "Ringkasan & Pembayaran - EduVera",
		"TenantID":       tenantID,
		"TenantName":     name,
		"Subdomain":      subdomain,
		"PackageName":    packageName,
		"Price":          price,
		"PriceFormatted": fmt.Sprintf("Rp %d", price),
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Gagal render template step5: "+err.Error(), http.StatusInternalServerError)
	}
}
