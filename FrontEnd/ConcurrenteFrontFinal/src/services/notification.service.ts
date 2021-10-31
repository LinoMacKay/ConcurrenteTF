import { Injectable } from '@angular/core';
import { MatSnackBar } from '@angular/material/snack-bar';
import { NotificationComponent } from 'src/app/shared/notification/notification.component';

@Injectable({
  providedIn: 'root',
})
export class NotificationService {
  constructor(private _matSnackbar: MatSnackBar) {}

  OpenSnackbar(message: string, type: string) {
    let icon;
    let snackbarStyle = '';

    if (type == 'error') {
      snackbarStyle = 'red-snackbar';
      icon = 'error';
    } else if (type == 'success') {
      snackbarStyle = 'green-snackbar';
      icon = 'done';
    } else {
      snackbarStyle = 'yellow-snackbar';
      icon = 'info';
    }

    this._matSnackbar.openFromComponent(NotificationComponent, {
      duration: 2000,
      panelClass: [snackbarStyle],
      data: {
        message: message,
        icon: icon,
      },
    });
  }
}
