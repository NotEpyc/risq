import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { MessageCircle, Send, HelpCircle } from "lucide-react";

interface QAItem {
  id: number;
  question: string;
  answer: string;
  timestamp: Date;
}

const QASection = () => {
  const [question, setQuestion] = useState("");
  const [email, setEmail] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);
  
  // Sample Q&A data - in a real app, this would come from a database
  const [qaItems] = useState<QAItem[]>([
    {
      id: 1,
      question: "How does your risk assessment platform work?",
      answer: "Our platform analyzes startup domains using advanced algorithms to evaluate various risk factors including market viability, financial stability, competitive landscape, and team expertise. We provide comprehensive reports with actionable insights.",
      timestamp: new Date('2024-01-15')
    },
    {
      id: 2,
      question: "What types of startups can I analyze?",
      answer: "You can analyze any startup with an online presence. Our platform works with companies across all industries and stages, from early-stage startups to established businesses looking to expand.",
      timestamp: new Date('2024-01-10')
    },
    {
      id: 3,
      question: "How accurate are your risk assessments?",
      answer: "Our assessments have an 85% success rate based on historical data. We use multiple data sources and machine learning algorithms to provide the most accurate risk evaluation possible.",
      timestamp: new Date('2024-01-08')
    }
  ]);

  const handleSubmitQuestion = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!question.trim()) return;

    setIsSubmitting(true);
    
    // Simulate API call
    await new Promise(resolve => setTimeout(resolve, 1000));
    
    console.log("Question submitted:", { question, email });
    
    // Reset form
    setQuestion("");
    setEmail("");
    setIsSubmitting(false);
    
    // In a real app, you would show a success toast here
    alert("Thank you for your question! We'll get back to you soon.");
  };

  return (
    <section className="py-20 bg-white">
      <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="text-center mb-12">
          <div className="inline-flex items-center gap-2 px-4 py-2 bg-blue-50 rounded-full border border-blue-200 text-blue-700 text-sm font-medium mb-4">
            <HelpCircle className="w-4 h-4" />
            Questions & Answers
          </div>
          <h2 className="text-3xl sm:text-4xl font-bold text-gray-900 mb-4">
            Have Questions About Our Platform?
          </h2>
          <p className="text-xl text-gray-600 max-w-2xl mx-auto">
            Get answers to common questions or ask your own. We're here to help you make informed decisions.
          </p>
        </div>

        <div className="grid gap-12 lg:grid-cols-2">
          {/* Question Submission Form */}
          <div className="order-2 lg:order-1">
            <Card className="shadow-lg">
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <MessageCircle className="w-5 h-5" />
                  Ask a Question
                </CardTitle>
              </CardHeader>
              <CardContent>
                <form onSubmit={handleSubmitQuestion} className="space-y-4">
                  <div>
                    <label htmlFor="email" className="block text-sm font-medium text-gray-700 mb-2">
                      Email Address (optional)
                    </label>
                    <Input
                      id="email"
                      type="email"
                      value={email}
                      onChange={(e) => setEmail(e.target.value)}
                      placeholder="your@email.com"
                      className="w-full"
                    />
                  </div>
                  
                  <div>
                    <label htmlFor="question" className="block text-sm font-medium text-gray-700 mb-2">
                      Your Question *
                    </label>
                    <Textarea
                      id="question"
                      value={question}
                      onChange={(e) => setQuestion(e.target.value)}
                      placeholder="What would you like to know about our risk assessment platform?"
                      className="w-full h-32"
                      required
                    />
                  </div>
                  
                  <Button 
                    type="submit" 
                    disabled={isSubmitting || !question.trim()}
                    className="w-full bg-gradient-to-r from-blue-600 to-purple-600 hover:from-blue-700 hover:to-purple-700"
                  >
                    {isSubmitting ? (
                      "Submitting..."
                    ) : (
                      <>
                        <Send className="w-4 h-4 mr-2" />
                        Submit Question
                      </>
                    )}
                  </Button>
                </form>
              </CardContent>
            </Card>
          </div>

          {/* Existing Q&A Display */}
          <div className="order-1 lg:order-2">
            <h3 className="text-2xl font-bold text-gray-900 mb-6">Frequently Asked Questions</h3>
            <div className="space-y-6">
              {qaItems.map((item) => (
                <Card key={item.id} className="shadow-sm hover:shadow-md transition-shadow">
                  <CardContent className="p-6">
                    <div className="mb-3">
                      <h4 className="font-semibold text-gray-900 mb-2 flex items-start gap-2">
                        <HelpCircle className="w-5 h-5 mt-0.5 text-blue-600 flex-shrink-0" />
                        {item.question}
                      </h4>
                    </div>
                    <p className="text-gray-600 leading-relaxed pl-7">
                      {item.answer}
                    </p>
                    <div className="mt-4 pt-3 border-t border-gray-100">
                      <span className="text-sm text-gray-400">
                        {item.timestamp.toLocaleDateString()}
                      </span>
                    </div>
                  </CardContent>
                </Card>
              ))}
            </div>
          </div>
        </div>

        <div className="mt-12 text-center">
          <p className="text-gray-600">
            Don't see your question here? Feel free to ask using the form above, and we'll get back to you as soon as possible.
          </p>
        </div>
      </div>
    </section>
  );
};

export default QASection;
