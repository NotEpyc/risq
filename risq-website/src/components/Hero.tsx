
import { Shield, TrendingUp, Users } from "lucide-react";
import SearchBar from "./SearchBar";

const Hero = () => {
  const handleSearch = (domain: string) => {
    console.log("Searching for startup domain:", domain);
    // Here you would typically integrate with your backend API
    // For now, we'll just log the search term
  };

  return (
    <section className="relative overflow-hidden bg-gradient-to-br from-slate-900 via-blue-900 to-indigo-900 text-white">
      <div className="absolute inset-0 opacity-20" style={{
        backgroundImage: `url("data:image/svg+xml,%3Csvg width='60' height='60' viewBox='0 0 60 60' xmlns='http://www.w3.org/2000/svg'%3E%3Cg fill='none' fill-rule='evenodd'%3E%3Cg fill='%239C92AC' fill-opacity='0.1'%3E%3Ccircle cx='30' cy='30' r='4'/%3E%3C/g%3E%3C/g%3E%3C/svg%3E")`
      }}></div>
      
      <div className="relative max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-20 lg:py-32">
        <div className="text-center space-y-8">
          <div className="inline-flex items-center gap-2 px-4 py-2 bg-blue-500/20 rounded-full border border-blue-400/30 text-blue-200 text-sm font-medium">
            <Shield className="w-4 h-4" />
            Smart Risk Assessment Platform
          </div>
          
          <h1 className="text-4xl sm:text-5xl lg:text-6xl font-bold leading-tight">
            Turn Startup Risks Into
            <span className="bg-gradient-to-r from-blue-400 via-purple-400 to-cyan-400 bg-clip-text text-transparent block mt-2">
              Strategic Advantages
            </span>
          </h1>
          
          <p className="text-xl text-slate-300 max-w-3xl mx-auto leading-relaxed">
            Empower your entrepreneurial journey with data-driven risk assessment. 
            Get personalized insights, prioritized action plans, and the confidence to make informed decisions from day one.
          </p>
          
          <div className="space-y-4">
            <SearchBar onSearch={handleSearch} />
            <p className="text-sm text-slate-400">
              Search and analyze any startup by their domain name
            </p>
          </div>
          
          <div className="flex justify-center items-center gap-8 pt-8 text-slate-400">
            <div className="flex items-center gap-2">
              <Users className="w-5 h-5" />
              <span>1,000+ Entrepreneurs</span>
            </div>
            <div className="flex items-center gap-2">
              <TrendingUp className="w-5 h-5" />
              <span>85% Success Rate</span>
            </div>
            <div className="flex items-center gap-2">
              <Shield className="w-5 h-5" />
              <span>Risk-Free Assessment</span>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
};

export default Hero;
