import { Component, OnInit } from '@angular/core';
import { FormControl, FormGroup, Validators } from '@angular/forms';
import { SymptomsService } from 'src/services/symptoms.service';

@Component({
  selector: 'app-prediction',
  templateUrl: './prediction.component.html',
  styleUrls: ['./prediction.component.scss'],
})
export class PredictionComponent implements OnInit {
  constructor(public symptomsService: SymptomsService) {}

  predictorForm = new FormGroup({});
  sintomas = this.symptomsService.symptoms;
  ngOnInit(): void {
    this.predictorForm = new FormGroup({
      name: new FormControl('', [Validators.required]),
    });
  }
  submitForm() {
    if (this.predictorForm.valid) {
      var toSend = {
        name: this.predictorForm.get('name')?.value,
        sintomas: this.symptomsService.getOnlySelected(),
      };
      console.log(toSend);
    }
  }
}
