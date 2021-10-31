import { Component, OnInit } from '@angular/core';
import { FormControl, FormGroup, Validators } from '@angular/forms';
import { MatDialog } from '@angular/material/dialog';
import { NotificationService } from 'src/services/notification.service';
import { SymptomsService } from 'src/services/symptoms.service';
import { ResultComponent } from '../result/result.component';

@Component({
  selector: 'app-prediction',
  templateUrl: './prediction.component.html',
  styleUrls: ['./prediction.component.scss'],
})
export class PredictionComponent implements OnInit {
  constructor(
    public symptomsService: SymptomsService,
    public dialog: MatDialog,
    private notificationService: NotificationService
  ) {}

  predictorForm = new FormGroup({});
  sintomas = this.symptomsService.symptoms;
  ngOnInit(): void {
    this.predictorForm = new FormGroup({
      name: new FormControl('', [Validators.required]),
    });
  }

  openDialog(toSend: any): void {
    const dialogRef = this.dialog.open(ResultComponent, {
      data: { toSend: toSend },
    });
  }

  submitForm() {
    if (
      true
      // this.predictorForm.valid &&
      //this.symptomsService.getOnlySelected().length > 0
    ) {
      var toSend = {
        name: this.predictorForm.get('name')?.value,
        sintomas: this.symptomsService.getOnlySelected(),
      };
      this.openDialog(toSend);
      console.log(toSend);
    } else {
      this.notificationService.OpenSnackbar(
        'Por favor complete correctamente el formulario',
        'warning'
      );
    }
  }
}
