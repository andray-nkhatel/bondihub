import { HttpInterceptorFn } from '@angular/common/http';
import { inject } from '@angular/core';
import { finalize } from 'rxjs';
import { NgxSpinnerService } from 'ngx-spinner';

export const loadingInterceptor: HttpInterceptorFn = (req, next) => {
  const spinner = inject(NgxSpinnerService);
  
  // Show spinner for non-GET requests or requests that might take time
  const shouldShowSpinner = req.method !== 'GET' || 
    req.url.includes('/upload') || 
    req.url.includes('/process') ||
    req.url.includes('/generate');

  if (shouldShowSpinner) {
    spinner.show();
  }

  return next(req).pipe(
    finalize(() => {
      if (shouldShowSpinner) {
        spinner.hide();
      }
    })
  );
};
