import { Component, Input, OnInit } from '@angular/core';
import { Symptom } from 'src/app/models/symptom';
import { SymptomsService } from 'src/services/symptoms.service';

@Component({
  selector: 'app-simtomas-card',
  templateUrl: './simtomas-card.component.html',
  styleUrls: ['./simtomas-card.component.scss'],
})
export class SimtomasCardComponent implements OnInit {
  constructor(private symptomsService: SymptomsService) {}
  @Input() sintoma = new Symptom('', '');
  ngOnInit(): void {}

  selectSymptom() {
    this.symptomsService.selectSymptom(this.sintoma);
  }
}
