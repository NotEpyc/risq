
import Hero from "@/components/Hero";
import Problem from "@/components/Problem";
import Features from "@/components/Features";
import Benefits from "@/components/Benefits";
import QASection from "@/components/QASection";
import CallToAction from "@/components/CallToAction";

const Index = () => {
  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-50 via-blue-50 to-indigo-50">
      <Hero />
      <Problem />
      <Features />
      <Benefits />
      <QASection />
      <CallToAction />
    </div>
  );
};

export default Index;
