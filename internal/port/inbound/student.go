package inbound_port

// StudentHttpPort defines the HTTP handlers for student endpoints
type StudentHttpPort interface {
	List(a any) error   // GET /api/v1/students
	Get(a any) error    // GET /api/v1/students/:id
	Create(a any) error // POST /api/v1/students
	Update(a any) error // PUT /api/v1/students/:id
	Delete(a any) error // DELETE /api/v1/students/:id
	Count(a any) error  // GET /api/v1/students/count
}
