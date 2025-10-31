export const environment = {
  production: true,
  apiUrl: 'https://api.bondihub.com/api/v1',
  appName: 'BondiHub',
  appVersion: '1.0.0',
  currency: 'ZMW',
  currencySymbol: 'K',
  supportEmail: 'support@bondihub.com',
  supportPhone: '+260 211 123 456',
  socialMedia: {
    facebook: 'https://facebook.com/bondihub',
    twitter: 'https://twitter.com/bondihub',
    instagram: 'https://instagram.com/bondihub',
    linkedin: 'https://linkedin.com/company/bondihub'
  },
  features: {
    enablePushNotifications: true,
    enableOfflineMode: true,
    enableAnalytics: true,
    enableErrorReporting: true
  },
  paymentMethods: {
    mtn: {
      enabled: true,
      name: 'MTN MoMo',
      icon: 'assets/icons/mtn-momo.svg'
    },
    airtel: {
      enabled: true,
      name: 'Airtel Money',
      icon: 'assets/icons/airtel-money.svg'
    },
    cash: {
      enabled: true,
      name: 'Cash',
      icon: 'assets/icons/cash.svg'
    },
    bank: {
      enabled: true,
      name: 'Bank Transfer',
      icon: 'assets/icons/bank.svg'
    }
  },
  map: {
    defaultCenter: {
      lat: -15.3875,
      lng: 28.3228
    },
    defaultZoom: 10,
    apiKey: 'your-production-google-maps-api-key'
  },
  cloudinary: {
    cloudName: 'your-production-cloud-name',
    uploadPreset: 'bondihub-uploads'
  }
};
