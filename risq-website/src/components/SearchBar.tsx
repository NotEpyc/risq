
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Search } from "lucide-react";
import { useState } from "react";

interface SearchBarProps {
  onSearch: (domain: string) => void;
  placeholder?: string;
}

const SearchBar = ({ onSearch, placeholder = "Enter startup domain (e.g., example.com)" }: SearchBarProps) => {
  const [searchTerm, setSearchTerm] = useState("");

  const handleSearch = () => {
    if (searchTerm.trim()) {
      onSearch(searchTerm.trim());
    }
  };

  const handleKeyPress = (e: React.KeyboardEvent) => {
    if (e.key === "Enter") {
      handleSearch();
    }
  };

  return (
    <div className="flex flex-col sm:flex-row gap-4 max-w-lg mx-auto">
      <Input
        type="text"
        value={searchTerm}
        onChange={(e) => setSearchTerm(e.target.value)}
        onKeyPress={handleKeyPress}
        placeholder={placeholder}
        className="flex-1 h-12 px-4 text-lg bg-white/10 border-white/20 text-white placeholder:text-white/60 focus:bg-white/20"
      />
      <Button 
        onClick={handleSearch}
        size="lg" 
        className="bg-gradient-to-r from-blue-600 to-purple-600 hover:from-blue-700 hover:to-purple-700 text-white px-8 py-3 text-lg font-semibold group transition-all duration-300 hover:scale-105"
      >
        <Search className="w-5 h-5 mr-2" />
        Search
      </Button>
    </div>
  );
};

export default SearchBar;
