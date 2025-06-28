
import { Star, Users, Zap } from "lucide-react";
import SearchBar from "./SearchBar";

const CallToAction = () => {
  const handleSearch = (domain: string) => {
    console.log("Searching for startup domain:", domain);
    // Here you would typically integrate with your backend API
    // For now, we'll just log the search term
  };

  return (
    <section className="py-20 bg-gradient-to-br from-slate-900 via-blue-900 to-indigo-900 text-white relative overflow-hidden">
      <div className="absolute inset-0 opacity-20" style={{
        backgroundImage: `url("data:image/svg+xml,%3Csvg width='60' height='60' viewBox='0 0 60 60' xmlns='http://www.w3.org/2000/svg'%3E%3Cg fill='none' fill-rule='evenodd'%3E%3Cg fill='%239C92AC' fill-opacity='0.1'%3E%3Cpath d='M30 30l15-15v30l-15-15zM30 30l-15 15h30l-15-15z'/%3E%3C/g%3E%3C/g%3E%3C/svg%3E")`
      }}></div>
      
      <div className="relative max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 text-center">
        <div className="space-y-8">
          <h2 className="text-3xl sm:text-4xl lg:text-5xl font-bold leading-tight">
            Ready to Analyze Your Next
            <span className="bg-gradient-to-r from-blue-400 via-purple-400 to-cyan-400 bg-clip-text text-transparent block mt-2">
              Investment Opportunity?
            </span>
          </h2>
          
          <p className="text-xl text-slate-300 max-w-2xl mx-auto leading-relaxed">
            Search for any startup by domain and get instant risk assessment insights to make informed investment decisions.
          </p>
          
          <div className="space-y-4">
            <SearchBar onSearch={handleSearch} placeholder="Search startup domain for analysis..." />
            <p className="text-sm text-slate-400">
              Enter any startup domain to begin your risk assessment
            </p>
          </div>
          
          <div className="flex justify-center items-center gap-8 pt-8 text-slate-400">
            <div className="flex items-center gap-2">
              <Star className="w-5 h-5 text-yellow-400" />
              <span>Instant Analysis</span>
            </div>
            <div className="flex items-center gap-2">
              <Zap className="w-5 h-5 text-green-400" />
              <span>Real-time Data</span>
            </div>
            <div className="flex items-center gap-2">
              <Users className="w-5 h-5 text-blue-400" />
              <span>Expert Insights</span>
            </div>
          </div>
          
          <div className="pt-8 border-t border-slate-700">
            <p className="text-slate-400 text-sm">
              Trusted by investors worldwide • No registration required • Start searching immediately
            </p>
          </div>
        </div>
      </div>
    </section>
  );
};

export default CallToAction;
