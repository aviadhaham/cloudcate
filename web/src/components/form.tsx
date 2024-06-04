import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Separator } from "@/components/ui/separator";
import { AllSearchResults } from "@/types/search-results";
import { LoaderCircle } from "lucide-react";
import { useState } from "react";

type Props = {
  searchQuery: string;
  setSearchQuery: (value: string) => void;
  typeQuery: string;
  setTypeQuery: (value: string) => void;
  subTypeQuery: string;
  setSubTypeQuery: (value: string) => void;
  onResults: (data: AllSearchResults[]) => void;
};

function isSearchQueryValid(query: string) {
  if (query === "") {
    return false;
  } else if (query.trim() === "") {
    return false;
  } else {
    return true;
  }
}

export default function Form(props: Props) {
  const [isLoading, setIsLoading] = useState(false);

  const handleTypeValueChange = (value: string) => {
    const types = value.split(":");
    const type = types[0];
    const subType = types[1] || "";
    props.setTypeQuery(type);
    props.setSubTypeQuery(subType);
  };

  const sendSearchRequest = async () => {
    // Clear the results when a new search is made
    props.onResults([]);
    setIsLoading(true);

    let url = `/api/search?resource_name=${props.searchQuery}&resource_type=${props.typeQuery}`;

    if (props.subTypeQuery !== "") {
      url += `&resource_subtype=${props.subTypeQuery}`;
    }

    const response = await fetch(url);

    if (!response.ok) {
      console.error(
        `Error fetching data: ${response.status} ${response.statusText}`
      );
      setIsLoading(false);
      return;
    }

    const contentType = response.headers.get("content-type");
    if (contentType && contentType.indexOf("application/json") !== -1) {
      try {
        setIsLoading(true);
        const data = await response.json();
        props.onResults(data.results);
      } catch (e) {
        console.error(e);
      } finally {
        setIsLoading(false);
      }
    } else {
      console.error("The response is not a valid JSON");
    }
  };

  return (
    <>
      <div className="flex items-center justify-center my-10 gap-x-6">
        <div className="w-[300px]">
          <label htmlFor="resource-name" className="block text-sm font-medium leading-6 text-gray-900">
            Resource Name (may be partial)
          </label>
          <div>
            <Input
              type="text"
              value={props.searchQuery}
              onChange={(e) => props.setSearchQuery(e.target.value)}
              name="resource-name"
              id="resource-name"
            />
          </div>
        </div>
        <div>
          <label htmlFor="resource-type" className="block text-sm font-medium leading-6 text-gray-900">
            Resource Type
          </label>
          <div>
            <Select
              value={`${props.typeQuery}${props.subTypeQuery ? `:${props.subTypeQuery}` : ""}`}
              onValueChange={handleTypeValueChange}
            >
              <SelectTrigger className="w-[360px]">
                <SelectValue placeholder="Select a type" />
              </SelectTrigger>
              <SelectContent>
                <SelectGroup>
                  <SelectItem value="s3">S3 Bucket</SelectItem>
                  <SelectItem value="dns">DNS (Hosted Zone or Record)</SelectItem>
                  <SelectItem value="loadbalancer">Load Balancer</SelectItem>
                  <SelectItem value="ec2">EC2 Instance (by ID, IP, DNS, or Tags)</SelectItem>
                  <SelectItem value="iam:key">IAM (Access Key)</SelectItem>
                  <SelectItem value="iam:user">IAM (User)</SelectItem>
                  <SelectItem value="elastic_ip">Elastic IP</SelectItem>
                  <SelectItem value="cloudfront">CloudFront Distribution (by ID or Domain name)</SelectItem>
                </SelectGroup>
              </SelectContent>
            </Select>
          </div>
        </div>
        <Button
          onClick={isSearchQueryValid(props.searchQuery) == true ? sendSearchRequest : undefined}
          className="self-end"
        >
          Search AWS
        </Button>
      </div>
      <div className="flex flex-col items-center">
        <Separator />
        {isLoading && (
          <div className="mt-8 animate-spin">
            <LoaderCircle className="w-10 h-10" />
          </div>
        )}
      </div>
    </>
  );
}
