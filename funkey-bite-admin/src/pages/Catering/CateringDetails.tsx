// src/pages/Catering/CateringDetails.tsx
import React, { useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import {
  ArrowLeft,
  Calendar,
  Users,
  Phone,
  Mail,
  DollarSign,
  MapPin,
  Clock,
  CheckCircle,
  XCircle,
  AlertCircle,
  Edit,
  MessageSquare,
  Printer,
  Download,
  ChevronDown,
  User,
  Package,
  FileText,
  Award,
  TrendingUp,
  MoreVertical,
  Plus
} from 'lucide-react';
import { Link } from 'react-router-dom';
import toast from 'react-hot-toast';
import { getCateringRequests, updateCateringStatus } from '../../api/adminApi';
import type { CateringRequest, CateringStatus } from '../../types';

const CateringDetails: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [isUpdating, setIsUpdating] = useState(false);

  // Fetch catering request details
  const { data: cateringRequest, isLoading, refetch } = useQuery({
    queryKey: ['catering-request', id],
    queryFn: async () => {
      if (!id) throw new Error('No catering request ID');
      try {
        const data = await getCateringRequests({ page: 1, limit: 1, id: parseInt(id) });
        const requests = data.data || [];
        const foundRequest = requests.find((r: CateringRequest) => r.id === parseInt(id));
        
        if (!foundRequest) {
          toast.error('Catering request not found');
          navigate('/catering');
          return null;
        }
        return foundRequest;
      } catch (error) {
        toast.error('Failed to load catering request details');
        navigate('/catering');
        return null;
      }
    },
    enabled: !!id,
  });

  const handleStatusUpdate = async (newStatus: CateringStatus) => {
    if (!cateringRequest || !id) return;
    
    setIsUpdating(true);
    try {
      await updateCateringStatus(parseInt(id), { status: newStatus });
      toast.success(`Catering request status updated to ${newStatus}`);
      refetch();
    } catch (error) {
      toast.error('Failed to update status');
    } finally {
      setIsUpdating(false);
    }
  };

  const getStatusInfo = (status: CateringStatus) => {
    switch (status) {
      case 'pending':
        return { 
          color: 'text-yellow-600 bg-yellow-50 dark:bg-yellow-900/20 border-yellow-200',
          icon: <AlertCircle className="h-5 w-5" />,
          text: 'Pending - Awaiting confirmation',
          action: 'Confirm Request'
        };
      case 'confirmed':
        return { 
          color: 'text-green-600 bg-green-50 dark:bg-green-900/20 border-green-200',
          icon: <CheckCircle className="h-5 w-5" />,
          text: 'Confirmed - Event scheduled',
          action: 'Mark as Completed'
        };
      case 'completed':
        return { 
          color: 'text-blue-600 bg-blue-50 dark:bg-blue-900/20 border-blue-200',
          icon: <Award className="h-5 w-5" />,
          text: 'Completed - Event fulfilled',
          action: 'Reopen Request'
        };
      case 'declined':
        return { 
          color: 'text-red-600 bg-red-50 dark:bg-red-900/20 border-red-200',
          icon: <XCircle className="h-5 w-5" />,
          text: 'Declined - Request declined',
          action: 'Reconsider Request'
        };
      default:
        return { 
          color: 'text-gray-600 bg-gray-50 dark:bg-gray-900/20 border-gray-200',
          icon: <AlertCircle className="h-5 w-5" />,
          text: 'Unknown status',
          action: 'Update Status'
        };
    }
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      weekday: 'long',
      year: 'numeric',
      month: 'long',
      day: 'numeric'
    });
  };

  const formatTime = (timeString?: string) => {
    if (!timeString) return 'All Day';
    return new Date(`2000-01-01T${timeString}`).toLocaleTimeString([], {
      hour: '2-digit',
      minute: '2-digit'
    });
  };

  const calculateEstimatedTotal = () => {
    if (!cateringRequest) return 0;
    
    if (cateringRequest.budget) {
      return cateringRequest.budget;
    }
    
    // Estimate based on guest count and package
    const perPersonRate = cateringRequest.package ? 
      (cateringRequest.package.includes('premium') ? 35 : 
       cateringRequest.package.includes('standard') ? 25 : 20) : 20;
    
    return cateringRequest.guestCount * perPersonRate;
  };

  const getPackageDetails = (packageName?: string) => {
    if (!packageName) return { name: 'Custom', description: 'Tailored catering package', pricePerPerson: 20 };
    
    const packages = {
      'basic': { name: 'Basic Package', description: 'Essential catering services', pricePerPerson: 20 },
      'standard': { name: 'Standard Package', description: 'Full catering services with setup', pricePerPerson: 25 },
      'premium': { name: 'Premium Package', description: 'Premium services with full coordination', pricePerPerson: 35 },
      'corporate': { name: 'Corporate Package', description: 'Business events and meetings', pricePerPerson: 30 },
      'wedding': { name: 'Wedding Package', description: 'Full wedding catering services', pricePerPerson: 40 },
    };
    
    return packages[packageName.toLowerCase() as keyof typeof packages] || 
           { name: packageName, description: 'Custom catering package', pricePerPerson: 25 };
  };

  if (isLoading) {
    return (
      <div className="space-y-6">
        <div className="flex items-center space-x-4">
          <div className="h-6 bg-gray-200 dark:bg-gray-700 rounded w-32 animate-pulse"></div>
        </div>
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          <div className="lg:col-span-2 space-y-4">
            {[...Array(3)].map((_, i) => (
              <div key={i} className="h-32 bg-gray-200 dark:bg-gray-700 rounded-xl animate-pulse"></div>
            ))}
          </div>
          <div className="space-y-4">
            {[...Array(2)].map((_, i) => (
              <div key={i} className="h-48 bg-gray-200 dark:bg-gray-700 rounded-xl animate-pulse"></div>
            ))}
          </div>
        </div>
      </div>
    );
  }

  if (!cateringRequest) {
    return (
      <div className="text-center py-12">
        <Calendar className="h-16 w-16 text-gray-400 mx-auto mb-4" />
        <h2 className="text-xl font-semibold text-gray-900 dark:text-white mb-2">
          Catering Request Not Found
        </h2>
        <p className="text-gray-500 dark:text-gray-400 mb-6">
          The catering request you're looking for doesn't exist.
        </p>
        <button
          onClick={() => navigate('/catering')}
          className="inline-flex items-center space-x-2 px-4 py-2 bg-primary-500 hover:bg-primary-600 text-white rounded-lg transition-colors"
        >
          <ArrowLeft className="h-4 w-4" />
          <span>Back to Catering</span>
        </button>
      </div>
    );
  }

  const statusInfo = getStatusInfo(cateringRequest.status);
  const packageDetails = getPackageDetails(cateringRequest.package);
  const estimatedTotal = calculateEstimatedTotal();
  const daysUntilEvent = Math.ceil(
    (new Date(cateringRequest.eventDate).getTime() - new Date().getTime()) / (1000 * 60 * 60 * 24)
  );

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
        <div className="flex items-center space-x-4">
          <button
            onClick={() => navigate('/catering')}
            className="p-2 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors"
          >
            <ArrowLeft className="h-5 w-5" />
          </button>
          <div>
            <h1 className="text-2xl font-bold text-gray-900 dark:text-white">
              {cateringRequest.eventName || 'Catering Request'}
            </h1>
            <p className="text-gray-600 dark:text-gray-400">
              Request ID: #{cateringRequest.id} • Submitted on {formatDate(cateringRequest.createdAt)}
            </p>
          </div>
        </div>
        
        <div className="flex items-center space-x-3">
          <button className="flex items-center space-x-2 px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors">
            <Printer className="h-4 w-4" />
            <span>Print</span>
          </button>
          <button className="flex items-center space-x-2 px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors">
            <Download className="h-4 w-4" />
            <span>Export</span>
          </button>
          <Link
            to={`/catering/${id}/edit`}
            className="flex items-center space-x-2 px-4 py-2 bg-primary-500 hover:bg-primary-600 text-white rounded-lg transition-colors"
          >
            <Edit className="h-4 w-4" />
            <span>Edit Request</span>
          </Link>
        </div>
      </div>

      {/* Status Banner */}
      <div className={`p-4 rounded-xl border ${statusInfo.color}`}>
        <div className="flex items-center justify-between">
          <div className="flex items-center space-x-3">
            {statusInfo.icon}
            <div>
              <p className="font-semibold">Current Status: {cateringRequest.status.toUpperCase()}</p>
              <p className="text-sm opacity-80">{statusInfo.text}</p>
              {daysUntilEvent >= 0 && cateringRequest.status === 'confirmed' && (
                <p className="text-sm mt-1">
                  Event in {daysUntilEvent} day{daysUntilEvent !== 1 ? 's' : ''}
                </p>
              )}
            </div>
          </div>
          
          <div className="relative group">
            <button 
              disabled={isUpdating}
              className="flex items-center space-x-2 px-4 py-2 bg-white dark:bg-gray-800 hover:bg-gray-50 dark:hover:bg-gray-700 rounded-lg transition-colors disabled:opacity-50"
            >
              <span>{statusInfo.action}</span>
              <ChevronDown className="h-4 w-4" />
            </button>
            <div className="absolute right-0 mt-2 w-48 bg-white dark:bg-gray-800 rounded-lg shadow-lg border border-gray-200 dark:border-gray-700 py-1 z-10 opacity-0 group-hover:opacity-100 transition-opacity pointer-events-none group-hover:pointer-events-auto">
              {cateringRequest.status !== 'confirmed' && (
                <button
                  onClick={() => handleStatusUpdate('confirmed')}
                  className="w-full text-left px-4 py-2 text-sm hover:bg-gray-100 dark:hover:bg-gray-700"
                >
                  Confirm Request
                </button>
              )}
              {cateringRequest.status !== 'completed' && (
                <button
                  onClick={() => handleStatusUpdate('completed')}
                  className="w-full text-left px-4 py-2 text-sm hover:bg-gray-100 dark:hover:bg-gray-700"
                >
                  Mark as Completed
                </button>
              )}
              {cateringRequest.status !== 'declined' && (
                <button
                  onClick={() => handleStatusUpdate('declined')}
                  className="w-full text-left px-4 py-2 text-sm text-red-600 dark:text-red-400 hover:bg-red-50 dark:hover:bg-red-900/20"
                >
                  Decline Request
                </button>
              )}
              {cateringRequest.status !== 'pending' && (
                <button
                  onClick={() => handleStatusUpdate('pending')}
                  className="w-full text-left px-4 py-2 text-sm hover:bg-gray-100 dark:hover:bg-gray-700"
                >
                  Reopen as Pending
                </button>
              )}
            </div>
          </div>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Left Column - Event Details */}
        <div className="lg:col-span-2 space-y-6">
          {/* Event Information */}
          <div className="bg-white dark:bg-gray-800 rounded-xl p-6 border border-gray-200 dark:border-gray-700">
            <h2 className="text-lg font-semibold text-gray-900 dark:text-white mb-6">
              Event Information
            </h2>
            
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div className="space-y-4">
                <div>
                  <label className="block text-sm font-medium text-gray-500 dark:text-gray-400 mb-1">
                    Event Name
                  </label>
                  <p className="font-medium text-gray-900 dark:text-white">
                    {cateringRequest.eventName || 'Not specified'}
                  </p>
                </div>
                
                <div>
                  <label className="block text-sm font-medium text-gray-500 dark:text-gray-400 mb-1">
                    Event Type
                  </label>
                  <p className="font-medium text-gray-900 dark:text-white">
                    {cateringRequest.eventType}
                  </p>
                </div>
                
                <div>
                  <label className="block text-sm font-medium text-gray-500 dark:text-gray-400 mb-1">
                    Package
                  </label>
                  <div className="flex items-center space-x-2">
                    <Package className="h-4 w-4 text-gray-400" />
                    <p className="font-medium text-gray-900 dark:text-white">
                      {packageDetails.name}
                    </p>
                  </div>
                  <p className="text-sm text-gray-500 dark:text-gray-400 mt-1">
                    {packageDetails.description}
                  </p>
                </div>
              </div>
              
              <div className="space-y-4">
                <div>
                  <label className="block text-sm font-medium text-gray-500 dark:text-gray-400 mb-1">
                    Date & Time
                  </label>
                  <div className="flex items-center space-x-2">
                    <Calendar className="h-4 w-4 text-gray-400" />
                    <div>
                      <p className="font-medium text-gray-900 dark:text-white">
                        {formatDate(cateringRequest.eventDate)}
                      </p>
                      <p className="text-sm text-gray-500 dark:text-gray-400">
                        {formatTime(cateringRequest.eventTime)}
                      </p>
                    </div>
                  </div>
                </div>
                
                <div>
                  <label className="block text-sm font-medium text-gray-500 dark:text-gray-400 mb-1">
                    Number of Guests
                  </label>
                  <div className="flex items-center space-x-2">
                    <Users className="h-4 w-4 text-gray-400" />
                    <p className="font-medium text-gray-900 dark:text-white">
                      {cateringRequest.guestCount} people
                    </p>
                  </div>
                </div>
                
                <div>
                  <label className="block text-sm font-medium text-gray-500 dark:text-gray-400 mb-1">
                    Budget
                  </label>
                  <div className="flex items-center space-x-2">
                    <DollarSign className="h-4 w-4 text-gray-400" />
                    <p className="font-medium text-gray-900 dark:text-white">
                      {cateringRequest.budget ? 
                        `$${cateringRequest.budget.toLocaleString()}` : 
                        'Not specified'}
                    </p>
                  </div>
                </div>
              </div>
            </div>
          </div>

          {/* Contact Information */}
          <div className="bg-white dark:bg-gray-800 rounded-xl p-6 border border-gray-200 dark:border-gray-700">
            <h2 className="text-lg font-semibold text-gray-900 dark:text-white mb-6">
              Contact Information
            </h2>
            
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div className="space-y-4">
                <div>
                  <label className="block text-sm font-medium text-gray-500 dark:text-gray-400 mb-1">
                    Contact Person
                  </label>
                  <div className="flex items-center space-x-2">
                    <User className="h-4 w-4 text-gray-400" />
                    <p className="font-medium text-gray-900 dark:text-white">
                      {cateringRequest.contactName}
                    </p>
                  </div>
                </div>
                
                <div>
                  <label className="block text-sm font-medium text-gray-500 dark:text-gray-400 mb-1">
                    Phone Number
                  </label>
                  <div className="flex items-center space-x-2">
                    <Phone className="h-4 w-4 text-gray-400" />
                    <p className="font-medium text-gray-900 dark:text-white">
                      {cateringRequest.contactPhone}
                    </p>
                  </div>
                </div>
              </div>
              
              <div className="space-y-4">
                <div>
                  <label className="block text-sm font-medium text-gray-500 dark:text-gray-400 mb-1">
                    Email Address
                  </label>
                  <div className="flex items-center space-x-2">
                    <Mail className="h-4 w-4 text-gray-400" />
                    <p className="font-medium text-gray-900 dark:text-white">
                      {cateringRequest.contactEmail || 'Not provided'}
                    </p>
                  </div>
                </div>
                
                {cateringRequest.userId && (
                  <div>
                    <label className="block text-sm font-medium text-gray-500 dark:text-gray-400 mb-1">
                      Customer Account
                    </label>
                    <Link
                      to={`/customers/${cateringRequest.userId}`}
                      className="inline-flex items-center space-x-2 text-primary-600 hover:text-primary-700 dark:text-primary-400 dark:hover:text-primary-300"
                    >
                      <span>View Customer Profile</span>
                      <TrendingUp className="h-4 w-4" />
                    </Link>
                  </div>
                )}
              </div>
            </div>
          </div>

          {/* Special Requests & Notes */}
          <div className="bg-white dark:bg-gray-800 rounded-xl p-6 border border-gray-200 dark:border-gray-700">
            <h2 className="text-lg font-semibold text-gray-900 dark:text-white mb-6">
              Special Requests & Notes
            </h2>
            
            <div className="space-y-6">
              {cateringRequest.specialRequests ? (
                <div className="p-4 bg-gray-50 dark:bg-gray-700 rounded-lg">
                  <div className="flex items-start space-x-3">
                    <FileText className="h-5 w-5 text-gray-400 mt-0.5" />
                    <div>
                      <p className="font-medium text-gray-900 dark:text-white mb-2">
                        Special Instructions
                      </p>
                      <p className="text-gray-700 dark:text-gray-300 whitespace-pre-line">
                        {cateringRequest.specialRequests}
                      </p>
                    </div>
                  </div>
                </div>
              ) : (
                <div className="text-center py-4">
                  <p className="text-gray-500 dark:text-gray-400">No special requests provided</p>
                </div>
              )}
              
              {/* Internal Notes Section */}
              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                  Internal Notes
                </label>
                <textarea
                  placeholder="Add internal notes about this catering request..."
                  rows={4}
                  className="w-full px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:bg-gray-700 dark:text-white"
                  defaultValue=""
                />
                <div className="flex justify-end mt-2">
                  <button className="px-4 py-2 text-sm bg-gray-100 dark:bg-gray-700 hover:bg-gray-200 dark:hover:bg-gray-600 rounded-lg transition-colors">
                    Save Note
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>

        {/* Right Column - Summary & Actions */}
        <div className="space-y-6">
          {/* Financial Summary */}
          <div className="bg-white dark:bg-gray-800 rounded-xl p-6 border border-gray-200 dark:border-gray-700">
            <h2 className="text-lg font-semibold text-gray-900 dark:text-white mb-6">
              Financial Summary
            </h2>
            
            <div className="space-y-4">
              <div className="flex justify-between items-center">
                <span className="text-gray-600 dark:text-gray-400">Package Rate</span>
                <span className="font-medium">${packageDetails.pricePerPerson}/person</span>
              </div>
              
              <div className="flex justify-between items-center">
                <span className="text-gray-600 dark:text-gray-400">Number of Guests</span>
                <span className="font-medium">{cateringRequest.guestCount}</span>
              </div>
              
              <div className="flex justify-between items-center">
                <span className="text-gray-600 dark:text-gray-400">Package Subtotal</span>
                <span className="font-medium">
                  ${(packageDetails.pricePerPerson * cateringRequest.guestCount).toLocaleString()}
                </span>
              </div>
              
              {cateringRequest.budget && (
                <div className="flex justify-between items-center">
                  <span className="text-gray-600 dark:text-gray-400">Client Budget</span>
                  <span className="font-medium">${cateringRequest.budget.toLocaleString()}</span>
                </div>
              )}
              
              <div className="border-t border-gray-200 dark:border-gray-700 pt-4">
                <div className="flex justify-between items-center">
                  <span className="text-lg font-semibold">Estimated Total</span>
                  <span className="text-2xl font-bold text-primary-600">
                    ${estimatedTotal.toLocaleString()}
                  </span>
                </div>
                <p className="text-sm text-gray-500 dark:text-gray-400 mt-2">
                  {cateringRequest.budget ? 
                    'Based on client budget' : 
                    'Estimated based on package and guest count'}
                </p>
              </div>
            </div>
          </div>

          {/* Quick Actions */}
          <div className="bg-white dark:bg-gray-800 rounded-xl p-6 border border-gray-200 dark:border-gray-700">
            <h2 className="text-lg font-semibold text-gray-900 dark:text-white mb-6">
              Quick Actions
            </h2>
            
            <div className="space-y-3">
              <button className="w-full flex items-center justify-between p-4 bg-gray-50 dark:bg-gray-700 hover:bg-gray-100 dark:hover:bg-gray-600 rounded-lg transition-colors">
                <div className="flex items-center space-x-3">
                  <MessageSquare className="h-5 w-5 text-gray-400" />
                  <span className="font-medium">Send Message</span>
                </div>
                <span className="text-sm text-gray-500">Email/SMS</span>
              </button>
              
              <button className="w-full flex items-center justify-between p-4 bg-gray-50 dark:bg-gray-700 hover:bg-gray-100 dark:hover:bg-gray-600 rounded-lg transition-colors">
                <div className="flex items-center space-x-3">
                  <Phone className="h-5 w-5 text-gray-400" />
                  <span className="font-medium">Call Client</span>
                </div>
                <span className="text-sm text-gray-500">{cateringRequest.contactPhone}</span>
              </button>
              
              <button className="w-full flex items-center justify-between p-4 bg-gray-50 dark:bg-gray-700 hover:bg-gray-100 dark:hover:bg-gray-600 rounded-lg transition-colors">
                <div className="flex items-center space-x-3">
                  <Calendar className="h-5 w-5 text-gray-400" />
                  <span className="font-medium">Add to Calendar</span>
                </div>
                <span className="text-sm text-gray-500">iCal/Google</span>
              </button>
              
              <button className="w-full flex items-center justify-between p-4 bg-gray-50 dark:bg-gray-700 hover:bg-gray-100 dark:hover:bg-gray-600 rounded-lg transition-colors">
                <div className="flex items-center space-x-3">
                  <FileText className="h-5 w-5 text-gray-400" />
                  <span className="font-medium">Generate Contract</span>
                </div>
                <span className="text-sm text-gray-500">PDF</span>
              </button>
            </div>
          </div>

          {/* Event Timeline */}
          <div className="bg-white dark:bg-gray-800 rounded-xl p-6 border border-gray-200 dark:border-gray-700">
            <h2 className="text-lg font-semibold text-gray-900 dark:text-white mb-6">
              Event Timeline
            </h2>
            
            <div className="space-y-4">
              <div className="flex items-start space-x-3">
                <div className="h-6 w-6 rounded-full bg-green-500 flex items-center justify-center mt-1">
                  <CheckCircle className="h-4 w-4 text-white" />
                </div>
                <div>
                  <p className="font-medium text-gray-900 dark:text-white">Request Submitted</p>
                  <p className="text-sm text-gray-500 dark:text-gray-400">
                    {formatDate(cateringRequest.createdAt)}
                  </p>
                </div>
              </div>
              
              {cateringRequest.status === 'confirmed' || cateringRequest.status === 'completed' ? (
                <div className="flex items-start space-x-3">
                  <div className="h-6 w-6 rounded-full bg-blue-500 flex items-center justify-center mt-1">
                    <CheckCircle className="h-4 w-4 text-white" />
                  </div>
                  <div>
                    <p className="font-medium text-gray-900 dark:text-white">Request Confirmed</p>
                    <p className="text-sm text-gray-500 dark:text-gray-400">
                      {cateringRequest.updatedAt ? formatDate(cateringRequest.updatedAt) : 'Date not recorded'}
                    </p>
                  </div>
                </div>
              ) : null}
              
              {cateringRequest.status === 'completed' ? (
                <div className="flex items-start space-x-3">
                  <div className="h-6 w-6 rounded-full bg-purple-500 flex items-center justify-center mt-1">
                    <Award className="h-4 w-4 text-white" />
                  </div>
                  <div>
                    <p className="font-medium text-gray-900 dark:text-white">Event Completed</p>
                    <p className="text-sm text-gray-500 dark:text-gray-400">
                      {cateringRequest.updatedAt ? formatDate(cateringRequest.updatedAt) : 'Date not recorded'}
                    </p>
                  </div>
                </div>
              ) : null}
              
              {cateringRequest.status === 'declined' ? (
                <div className="flex items-start space-x-3">
                  <div className="h-6 w-6 rounded-full bg-red-500 flex items-center justify-center mt-1">
                    <XCircle className="h-4 w-4 text-white" />
                  </div>
                  <div>
                    <p className="font-medium text-gray-900 dark:text-white">Request Declined</p>
                    <p className="text-sm text-gray-500 dark:text-gray-400">
                      {cateringRequest.updatedAt ? formatDate(cateringRequest.updatedAt) : 'Date not recorded'}
                    </p>
                  </div>
                </div>
              ) : null}
              
              <div className="flex items-start space-x-3">
                <div className="h-6 w-6 rounded-full bg-primary-500 flex items-center justify-center mt-1">
                  <Calendar className="h-4 w-4 text-white" />
                </div>
                <div>
                  <p className="font-medium text-gray-900 dark:text-white">Event Date</p>
                  <p className="text-sm text-gray-500 dark:text-gray-400">
                    {formatDate(cateringRequest.eventDate)} • {formatTime(cateringRequest.eventTime)}
                  </p>
                </div>
              </div>
            </div>
          </div>

          {/* Related Documents */}
          <div className="bg-white dark:bg-gray-800 rounded-xl p-6 border border-gray-200 dark:border-gray-700">
            <div className="flex items-center justify-between mb-6">
              <h2 className="text-lg font-semibold text-gray-900 dark:text-white">
                Related Documents
              </h2>
              <button className="text-sm text-primary-600 hover:text-primary-700 dark:text-primary-400 dark:hover:text-primary-300">
                Upload
              </button>
            </div>
            
            <div className="space-y-3">
              <div className="flex items-center justify-between p-3 bg-gray-50 dark:bg-gray-700 rounded-lg">
                <div className="flex items-center space-x-3">
                  <FileText className="h-5 w-5 text-gray-400" />
                  <div>
                    <p className="font-medium text-gray-900 dark:text-white">Event Contract</p>
                    <p className="text-xs text-gray-500 dark:text-gray-400">PDF • 2.4 MB</p>
                  </div>
                </div>
                <button className="text-gray-400 hover:text-gray-600 dark:hover:text-gray-300">
                  <Download className="h-4 w-4" />
                </button>
              </div>
              
              <div className="flex items-center justify-between p-3 bg-gray-50 dark:bg-gray-700 rounded-lg">
                <div className="flex items-center space-x-3">
                  <FileText className="h-5 w-5 text-gray-400" />
                  <div>
                    <p className="font-medium text-gray-900 dark:text-white">Menu Proposal</p>
                    <p className="text-xs text-gray-500 dark:text-gray-400">PDF • 1.8 MB</p>
                  </div>
                </div>
                <button className="text-gray-400 hover:text-gray-600 dark:hover:text-gray-300">
                  <Download className="h-4 w-4" />
                </button>
              </div>
              
              <div className="text-center py-4 border border-dashed border-gray-300 dark:border-gray-600 rounded-lg hover:border-primary-500 dark:hover:border-primary-500 transition-colors cursor-pointer">
                <div className="flex flex-col items-center space-y-2">
                  <div className="h-10 w-10 rounded-full bg-gray-100 dark:bg-gray-700 flex items-center justify-center">
                    <Plus className="h-5 w-5 text-gray-400" />
                  </div>
                  <p className="text-sm text-gray-600 dark:text-gray-400">Upload Document</p>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default CateringDetails;