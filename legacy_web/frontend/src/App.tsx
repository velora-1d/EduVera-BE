import { Button } from "@/components/ui/button";
import { useState } from "react";

function App() {
  const [count, setCount] = useState(0);

  return (
    <div className="flex flex-col items-center justify-center min-h-screen bg-background text-foreground">
      <h1 className="text-4xl font-bold mb-8">EduVkera Frontend</h1>
      <div className="card">
        <Button
          variant="default"
          onClick={() => setCount((count) => count + 1)}
          className="text-lg"
        >
          count is {count}
        </Button>
        <p className="mt-4 text-gray-500">React + Tailwind + JokoUI</p>
      </div>
    </div>
  );
}

export default App;
