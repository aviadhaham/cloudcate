import Form from "@/components/form";
import Header from "@/components/header";
import ResultsTable from "@/components/results-table";
import { AllSearchResults } from "@/types/search-results";
import { useState } from "react";

export default function HomePage() {
  const [results, setResults] = useState<AllSearchResults[]>([]);
  const handleResults = (data: AllSearchResults[]) => {
    setResults(data);
  };

  return (
    <div className="flex flex-col h-screen">
      <div className="flex flex-col w-full max-w-5xl mx-auto mt-5">
        <Header />
        <Form onResults={handleResults} />
        {results && results.length !== 0 && (
          <ResultsTable results={results} />
        )}
        {results === null && (
          <span className="mt-10 text-center">No results found.</span>
        )}
      </div>
    </div>
  );
}
