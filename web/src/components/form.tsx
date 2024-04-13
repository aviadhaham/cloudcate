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
import React, { useState } from "react";

type Props = {
  onResults: (data: AllSearchResults[]) => void;
};

export default function Form(props: Props) {
  const [isLoading, setIsLoading] = useState(false);
  const [searchTerm, setSearchTerm] = useState("");
  const [typeValue, setTypeValue] = useState("");

  const handleSearchTermChange = (
    event: React.ChangeEvent<HTMLInputElement>
  ) => {
    setSearchTerm(event.target.value);
  };

  const handleTypeValueChange = (value: string) => {
    setTypeValue(value);
  };

  const sendSearchRequest = async () => {
    // Clear the results when a new search is made
    props.onResults([]);
    setIsLoading(true);

    const response = await fetch(
      `/api/search?resource_name=${searchTerm}&resource_type=${typeValue}`
    );

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
          <label
            htmlFor="resource-name"
            className="block text-sm font-medium leading-6 text-gray-900"
          >
            Resource Name (may be partial)
          </label>
          <div>
            <Input
              type="text"
              value={searchTerm}
              onChange={handleSearchTermChange}
              name="resource-name"
              id="resource-name"
            />
          </div>
        </div>
        <div>
          <label
            htmlFor="resource-type"
            className="block text-sm font-medium leading-6 text-gray-900"
          >
            Resource Type
          </label>
          <div>
            <Select onValueChange={handleTypeValueChange}>
              <SelectTrigger className="w-[360px]">
                <SelectValue placeholder="Select a type" />
              </SelectTrigger>
              <SelectContent>
                <SelectGroup>
                  <SelectItem value="s3">S3 Bucket</SelectItem>
                  <SelectItem value="dns">
                    DNS (Hosted Zone or Record)
                  </SelectItem>
                  <SelectItem value="loadbalancer">Load Balancer</SelectItem>
                  <SelectItem value="ec2">
                    EC2 Instance (by IP, DNS, or Tags)
                  </SelectItem>
                  <SelectItem value="iam">IAM (Access Key)</SelectItem>
                  <SelectItem value="elastic_ip">Elastic IP</SelectItem>
                  <SelectItem value="cloudfront">
                    CloudFront Distribution (by ID or Domain name)
                  </SelectItem>
                </SelectGroup>
              </SelectContent>
            </Select>
          </div>
        </div>
        <Button onClick={sendSearchRequest} className="self-end">
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
