
import { CheckCircle, TrendingUp, Clock, Shield } from "lucide-react";

const Benefits = () => {
  const outcomes = [
    {
      icon: CheckCircle,
      title: "Informed Decision-Making",
      description: "Gain a clear, prioritized risk profile that transforms uncertainty into strategic advantage.",
      stat: "3x Better",
      statLabel: "Decision Quality"
    },
    {
      icon: Shield,
      title: "Proactive Risk Mitigation",
      description: "Address high-priority risks early with actionable suggestions that prevent critical failures.",
      stat: "75% Fewer",
      statLabel: "Critical Issues"
    },
    {
      icon: Clock,
      title: "Operational Efficiency",
      description: "Enable rapid feedback and iteration cycles that accelerate your path to product-market fit.",
      stat: "50% Faster",
      statLabel: "Time to Market"
    },
    {
      icon: TrendingUp,
      title: "Increased Success Rate",
      description: "Join the 15% of startups that succeed by making data-driven decisions from day one.",
      stat: "5x Higher",
      statLabel: "Success Probability"
    }
  ];

  return (
    <section className="py-20 bg-white">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="text-center mb-16">
          <h2 className="text-3xl sm:text-4xl font-bold text-slate-900 mb-6">
            Expected Outcomes That Drive Real Results
          </h2>
          <p className="text-xl text-slate-600 max-w-3xl mx-auto">
            Don't just take our word for it. See the measurable impact our platform has on entrepreneurial success.
          </p>
        </div>
        
        <div className="grid md:grid-cols-2 gap-8">
          {outcomes.map((outcome, index) => (
            <div 
              key={index} 
              className="group flex gap-6 p-8 bg-gradient-to-br from-slate-50 to-blue-50 rounded-2xl border border-slate-200 hover:border-blue-200 hover:shadow-lg transition-all duration-300"
            >
              <div className="flex-shrink-0">
                <div className="w-14 h-14 bg-gradient-to-br from-blue-500 to-purple-500 rounded-2xl flex items-center justify-center group-hover:scale-110 transition-transform duration-300">
                  <outcome.icon className="w-7 h-7 text-white" />
                </div>
              </div>
              
              <div className="flex-1">
                <h3 className="text-xl font-semibold text-slate-900 mb-3">{outcome.title}</h3>
                <p className="text-slate-600 mb-4 leading-relaxed">{outcome.description}</p>
                
                <div className="flex items-center gap-3">
                  <div className="text-2xl font-bold bg-gradient-to-r from-blue-600 to-purple-600 bg-clip-text text-transparent">
                    {outcome.stat}
                  </div>
                  <div className="text-sm text-slate-500 font-medium">
                    {outcome.statLabel}
                  </div>
                </div>
              </div>
            </div>
          ))}
        </div>
        
        <div className="mt-16 text-center">
          <div className="inline-flex items-center gap-4 px-8 py-4 bg-gradient-to-r from-green-50 to-emerald-50 rounded-xl border border-green-200">
            <TrendingUp className="w-6 h-6 text-green-600" />
            <span className="text-green-800 font-medium text-lg">
              Join 1,000+ entrepreneurs who've transformed their startup journey
            </span>
          </div>
        </div>
      </div>
    </section>
  );
};

export default Benefits;
