
import { Brain, Target, BarChart3, Lightbulb, Zap, CheckCircle } from "lucide-react";

const Features = () => {
  const features = [
    {
      icon: Brain,
      title: "AI-Powered Risk Analysis",
      description: "Advanced algorithms analyze your startup across multiple dimensions including business model, market dynamics, and financial exposure.",
      gradient: "from-blue-500 to-cyan-500"
    },
    {
      icon: Target,
      title: "Personalized Risk Profiling",
      description: "Get a customized risk assessment based on your specific industry, location, experience level, and business model.",
      gradient: "from-purple-500 to-pink-500"
    },
    {
      icon: BarChart3,
      title: "Priority-Based Insights",
      description: "Automatically prioritize risks by impact and likelihood, so you focus your limited resources on what matters most.",
      gradient: "from-green-500 to-emerald-500"
    },
    {
      icon: Lightbulb,
      title: "Actionable Recommendations",
      description: "Receive specific, implementable strategies to mitigate each identified risk, with step-by-step guidance.",
      gradient: "from-orange-500 to-red-500"
    },
    {
      icon: Zap,
      title: "Real-Time Updates",
      description: "Your risk profile evolves as your startup grows, with continuous monitoring and updated recommendations.",
      gradient: "from-indigo-500 to-blue-500"
    },
    {
      icon: CheckCircle,
      title: "Success Tracking",
      description: "Monitor your progress in addressing risks and see how your overall risk score improves over time.",
      gradient: "from-teal-500 to-green-500"
    }
  ];

  return (
    <section className="py-20 bg-gradient-to-br from-slate-50 to-blue-50">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="text-center mb-16">
          <h2 className="text-3xl sm:text-4xl font-bold text-slate-900 mb-6">
            Smart Features That Transform Risk Into Opportunity
          </h2>
          <p className="text-xl text-slate-600 max-w-3xl mx-auto">
            Our platform combines cutting-edge technology with entrepreneurial expertise to deliver insights that actually matter.
          </p>
        </div>
        
        <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-8">
          {features.map((feature, index) => (
            <div 
              key={index} 
              className="group relative p-8 bg-white rounded-2xl border border-slate-200 hover:border-slate-300 hover:shadow-xl transition-all duration-500 hover:-translate-y-2"
            >
              <div className={`w-14 h-14 bg-gradient-to-br ${feature.gradient} rounded-2xl flex items-center justify-center mb-6 group-hover:scale-110 transition-transform duration-300`}>
                <feature.icon className="w-7 h-7 text-white" />
              </div>
              <h3 className="text-xl font-semibold text-slate-900 mb-4">{feature.title}</h3>
              <p className="text-slate-600 leading-relaxed">{feature.description}</p>
              
              <div className={`absolute inset-0 bg-gradient-to-br ${feature.gradient} opacity-0 group-hover:opacity-5 rounded-2xl transition-opacity duration-300`}></div>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
};

export default Features;
