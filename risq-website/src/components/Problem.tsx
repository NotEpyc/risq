
import { AlertTriangle, TrendingDown, DollarSign, Clock } from "lucide-react";

const Problem = () => {
  const problems = [
    {
      icon: AlertTriangle,
      title: "Blind Spot Risks",
      description: "First-time entrepreneurs often miss critical risks that experienced founders spot immediately, leading to costly oversights."
    },
    {
      icon: TrendingDown,
      title: "Poor Decision-Making",
      description: "Without data-driven insights, entrepreneurs make gut decisions that can drain resources and derail their startup journey."
    },
    {
      icon: DollarSign,
      title: "Resource Depletion",
      description: "Ineffective risk prioritization leads to wasted time, money, and energy on low-impact activities while critical issues go unaddressed."
    },
    {
      icon: Clock,
      title: "Increased Failure Rate",
      description: "Statistics show that 90% of startups fail, often due to preventable risks that could have been identified and mitigated early."
    }
  ];

  return (
    <section className="py-20 bg-white">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="text-center mb-16">
          <h2 className="text-3xl sm:text-4xl font-bold text-slate-900 mb-6">
            The Hidden Challenges First-Time Entrepreneurs Face
          </h2>
          <p className="text-xl text-slate-600 max-w-3xl mx-auto">
            Starting a business is already challenging enough. Don't let preventable risks become the reason your startup doesn't succeed.
          </p>
        </div>
        
        <div className="grid md:grid-cols-2 lg:grid-cols-4 gap-8">
          {problems.map((problem, index) => (
            <div 
              key={index} 
              className="group p-6 bg-gradient-to-br from-slate-50 to-slate-100 rounded-2xl border border-slate-200 hover:border-red-200 hover:shadow-lg transition-all duration-300 hover:-translate-y-1"
            >
              <div className="w-12 h-12 bg-gradient-to-br from-red-100 to-orange-100 rounded-xl flex items-center justify-center mb-4 group-hover:scale-110 transition-transform duration-300">
                <problem.icon className="w-6 h-6 text-red-600" />
              </div>
              <h3 className="text-lg font-semibold text-slate-900 mb-2">{problem.title}</h3>
              <p className="text-slate-600 text-sm leading-relaxed">{problem.description}</p>
            </div>
          ))}
        </div>
        
        <div className="mt-16 text-center">
          <div className="inline-flex items-center gap-4 px-6 py-4 bg-gradient-to-r from-red-50 to-orange-50 rounded-xl border border-red-200">
            <AlertTriangle className="w-6 h-6 text-red-600" />
            <span className="text-red-800 font-medium">
              Don't become another statistic. Take control of your startup's future.
            </span>
          </div>
        </div>
      </div>
    </section>
  );
};

export default Problem;
