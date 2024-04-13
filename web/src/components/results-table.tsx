import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { AllSearchResults } from "@/types/search-results";

export default function ResultsTable({
  results,
}: {
  results: AllSearchResults[];
}) {
  if (results.length === 0) {
    return;
  }
  return (
    <Table className="mt-10">
      <TableHeader>
        <TableRow>
          {Object.keys(results[0]).map((key) => (
            <TableHead>{key}</TableHead>
          ))}
        </TableRow>
      </TableHeader>
      <TableBody>
        {results.map((result, rowIndex) => (
          <TableRow key={rowIndex}>
            {Object.keys(result).map((key, keyIndex) => (
              <TableCell key={keyIndex}>
                {result[key as keyof AllSearchResults]}
              </TableCell>
            ))}
          </TableRow>
        ))}
      </TableBody>
    </Table>
  );
}
