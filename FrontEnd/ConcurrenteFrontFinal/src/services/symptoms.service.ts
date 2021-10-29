import { Injectable } from '@angular/core';
import { Symptom } from 'src/app/models/symptom';

@Injectable({
  providedIn: 'root',
})
export class SymptomsService {
  public symptoms = [
    new Symptom('Gripe', ''),
    new Symptom('Tos', ''),
    new Symptom('Fiebre', ''),
    new Symptom('Dolor de garganta', ''),
    new Symptom('Dolor de cabeza', ''),
    new Symptom('Dolor de garganta', ''),
    new Symptom('Dolor de garganta', ''),
  ];

  constructor() {}

  selectSymptom(symptom: Symptom) {
    var index = this.symptoms.indexOf(symptom);
    this.symptoms[index].isSelected = !this.symptoms[index].isSelected;
  }

  GetSymtoms() {
    return this.symptoms;
  }

  getOnlySelected() {
    return this.symptoms.filter((e) => e.isSelected == true);
  }
}
