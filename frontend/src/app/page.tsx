import Link from 'next/link';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import {
  Users,
  FileText,
  AlertTriangle,
  XCircle,
  Building2,
  Calendar
} from 'lucide-react';

export default function DashboardPage() {
  // Mock data - will be replaced with real API calls
  const metrics = {
    totalEmployees: 248,
    activeDocuments: 212,
    expiringSoon: 28,
    expired: 8,
  };

  const alerts = [
    {
      id: '1',
      employeeName: 'John Doe',
      companyName: 'ABC Construction Co.',
      documentType: 'Visa',
      expiryDate: '2026-03-15',
      daysLeft: 30,
      status: 'warning' as const,
    },
    {
      id: '2',
      employeeName: 'Jane Smith',
      companyName: 'ABC Construction Co.',
      documentType: 'Passport',
      expiryDate: '2026-02-19',
      daysLeft: 7,
      status: 'urgent' as const,
    },
    {
      id: '3',
      employeeName: 'Mike Johnson',
      companyName: 'XYZ Engineering',
      documentType: 'Emirates ID',
      expiryDate: '2026-02-10',
      daysLeft: -2,
      status: 'expired' as const,
    },
  ];

  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-50 to-slate-100 p-4 sm:p-6 md:p-8">
      <div className="max-w-7xl mx-auto space-y-6 md:space-y-8">
        {/* Header - Responsive */}
        <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
          <div>
            <h1 className="text-3xl sm:text-4xl font-bold text-slate-900">Dashboard</h1>
            <p className="text-slate-600 mt-1 sm:mt-2 text-sm sm:text-base">Manpower Management System</p>
          </div>
          <div className="flex gap-3">
            <Link href="/employees" className="w-full sm:w-auto">
              <Button size="lg" className="w-full sm:w-auto">
                <Users className="h-4 w-4 mr-2" />
                View Employees
              </Button>
            </Link>
          </div>
        </div>

        {/* Metrics Cards - Fully Responsive Grid */}
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4 md:gap-6">
          <Card className="border-l-4 border-l-blue-500 shadow-lg hover:shadow-xl transition-shadow">
            <CardHeader className="flex flex-row items-center justify-between pb-2 space-y-0">
              <CardTitle className="text-sm font-medium text-slate-600">
                Total Employees
              </CardTitle>
              <Users className="h-5 w-5 text-blue-500" />
            </CardHeader>
            <CardContent>
              <div className="text-3xl font-bold text-slate-900">{metrics.totalEmployees}</div>
              <p className="text-xs text-slate-500 mt-1">Active workforce</p>
            </CardContent>
          </Card>

          <Card className="border-l-4 border-l-green-500 shadow-lg hover:shadow-xl transition-shadow">
            <CardHeader className="flex flex-row items-center justify-between pb-2 space-y-0">
              <CardTitle className="text-sm font-medium text-slate-600">
                Active Documents
              </CardTitle>
              <FileText className="h-5 w-5 text-green-500" />
            </CardHeader>
            <CardContent>
              <div className="text-3xl font-bold text-slate-900">{metrics.activeDocuments}</div>
              <p className="text-xs text-slate-500 mt-1">Valid documents</p>
            </CardContent>
          </Card>

          <Card className="border-l-4 border-l-yellow-500 shadow-lg hover:shadow-xl transition-shadow">
            <CardHeader className="flex flex-row items-center justify-between pb-2 space-y-0">
              <CardTitle className="text-sm font-medium text-slate-600">
                Expiring Soon
              </CardTitle>
              <AlertTriangle className="h-5 w-5 text-yellow-500" />
            </CardHeader>
            <CardContent>
              <div className="text-3xl font-bold text-yellow-700">{metrics.expiringSoon}</div>
              <p className="text-xs text-slate-500 mt-1">Within 30 days</p>
            </CardContent>
          </Card>

          <Card className="border-l-4 border-l-red-500 shadow-lg hover:shadow-xl transition-shadow">
            <CardHeader className="flex flex-row items-center justify-between pb-2 space-y-0">
              <CardTitle className="text-sm font-medium text-slate-600">
                Expired
              </CardTitle>
              <XCircle className="h-5 w-5 text-red-500" />
            </CardHeader>
            <CardContent>
              <div className="text-3xl font-bold text-red-700">{metrics.expired}</div>
              <p className="text-xs text-slate-500 mt-1">Requires renewal</p>
            </CardContent>
          </Card>
        </div>

        {/* Expiry Alerts - Mobile Optimized */}
        <Card className="shadow-lg">
          <CardHeader>
            <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3">
              <div>
                <CardTitle className="text-xl sm:text-2xl">Critical Expiry Alerts</CardTitle>
                <CardDescription className="mt-1 text-sm">
                  Documents expiring soon or already expired
                </CardDescription>
              </div>
              <Button variant="outline" size="sm" className="w-full sm:w-auto">
                <Calendar className="h-4 w-4 mr-2" />
                Filter
              </Button>
            </div>
          </CardHeader>
          <CardContent className="px-4 sm:px-6">
            <div className="space-y-3 sm:space-y-4">
              {alerts.map((alert) => (
                <div
                  key={alert.id}
                  className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3 sm:gap-4 p-4 rounded-lg border bg-card hover:shadow-md transition-shadow"
                >
                  {/* Left Section */}
                  <div className="flex-1 space-y-2">
                    <div className="flex flex-col sm:flex-row sm:items-center gap-2">
                      <h3 className="font-semibold text-base sm:text-lg text-slate-900">
                        {alert.employeeName}
                      </h3>
                      <Badge
                        variant="outline"
                        className={`w-fit ${alert.status === 'expired'
                            ? 'bg-red-100 text-red-800 border-red-200'
                            : alert.status === 'urgent'
                              ? 'bg-orange-100 text-orange-800 border-orange-200'
                              : 'bg-yellow-100 text-yellow-800 border-yellow-200'
                          }`}
                      >
                        {alert.status === 'expired'
                          ? 'Expired'
                          : alert.status === 'urgent'
                            ? 'Urgent'
                            : 'Warning'}
                      </Badge>
                    </div>
                    <div className="flex flex-col sm:flex-row sm:items-center gap-1 sm:gap-4 text-xs sm:text-sm text-slate-600">
                      <span className="flex items-center gap-1">
                        <Building2 className="h-3 w-3 flex-shrink-0" />
                        <span className="truncate">{alert.companyName}</span>
                      </span>
                      <span className="hidden sm:inline">•</span>
                      <span className="font-medium">{alert.documentType}</span>
                      <span className="hidden sm:inline">•</span>
                      <span className="text-xs">Expires: {alert.expiryDate}</span>
                    </div>
                  </div>

                  {/* Right Section - Days Counter */}
                  <div className="flex items-center justify-between sm:justify-end gap-4">
                    <div className="text-center sm:text-right">
                      <div
                        className={`text-2xl sm:text-3xl font-bold ${alert.daysLeft < 0
                            ? 'text-red-700'
                            : alert.daysLeft <= 7
                              ? 'text-orange-700'
                              : 'text-yellow-700'
                          }`}
                      >
                        {alert.daysLeft < 0 ? `${Math.abs(alert.daysLeft)}` : alert.daysLeft}
                      </div>
                      <div className="text-xs text-slate-500 whitespace-nowrap">
                        {alert.daysLeft < 0 ? 'days overdue' : 'days left'}
                      </div>
                    </div>
                    <Button variant="outline" size="sm">
                      View
                    </Button>
                  </div>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>

        {/* Company Summary - Mobile Optimized */}
        <Card className="shadow-lg">
          <CardHeader>
            <CardTitle className="text-xl sm:text-2xl">Company Summary</CardTitle>
            <CardDescription className="text-sm">Employee distribution across companies</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="space-y-3 sm:space-y-4">
              <div className="flex items-center justify-between p-3 rounded-lg bg-slate-50">
                <div className="flex items-center gap-2 sm:gap-3 flex-1 min-w-0">
                  <Building2 className="h-5 w-5 text-slate-600 flex-shrink-0" />
                  <span className="font-medium text-sm sm:text-base truncate">ABC Construction Co.</span>
                </div>
                <Badge variant="secondary" className="flex-shrink-0 ml-2">120 employees</Badge>
              </div>
              <div className="flex items-center justify-between p-3 rounded-lg bg-slate-50">
                <div className="flex items-center gap-2 sm:gap-3 flex-1 min-w-0">
                  <Building2 className="h-5 w-5 text-slate-600 flex-shrink-0" />
                  <span className="font-medium text-sm sm:text-base truncate">XYZ Engineering</span>
                </div>
                <Badge variant="secondary" className="flex-shrink-0 ml-2">85 employees</Badge>
              </div>
              <div className="flex items-center justify-between p-3 rounded-lg bg-slate-50">
                <div className="flex items-center gap-2 sm:gap-3 flex-1 min-w-0">
                  <Building2 className="h-5 w-5 text-slate-600 flex-shrink-0" />
                  <span className="font-medium text-sm sm:text-base truncate">Others</span>
                </div>
                <Badge variant="secondary" className="flex-shrink-0 ml-2">43 employees</Badge>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
