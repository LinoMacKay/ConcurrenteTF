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
  @Input() sintoma = new Symptom('', '', '');
  ngOnInit(): void {}
  background = 'linear-gradient(rgba(0, 0, 0, 0.5), rgba(0, 0, 0, 0.5))';
  image = this.sintoma.image;
  selectSymptom() {
    this.symptomsService.selectSymptom(this.sintoma);
  }
  myStyle(): object {
    return {
      'background-image': `${this.background + ',url(' + this.sintoma.image})`,
      'background-size': 'cover',
      'background-position': 'center',
    };
  }
}
