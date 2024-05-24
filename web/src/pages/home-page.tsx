import Form from "@/components/form";
import Header from "@/components/header";
import ResultsTable from "@/components/results-table";
import { AllSearchResults } from "@/types/search-results";
import { useEffect, useState } from "react";
import { useSearchParams } from "react-router-dom";

export default function HomePage() {
  const [searchParams, setSearchParams] = useSearchParams();
  const [searchQuery, setSearchQuery] = useState(searchParams.get("q") || "");
  const [results, setResults] = useState<AllSearchResults[]>([]);
  const handleResults = (data: AllSearchResults[]) => {
    setResults(data);
  };

  useEffect(() => {
    setSearchParams({ q: searchQuery });
  }, [setSearchParams, searchQuery]);

  return (
    <div className="flex flex-col h-screen">
      <div className="flex flex-col w-full max-w-5xl mx-auto mt-5">
        <Header />
        <Form
          searchQuery={searchQuery}
          setSearchQuery={setSearchQuery}
          onResults={handleResults}
        />
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
