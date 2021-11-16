import { Component, Inject, OnInit } from '@angular/core';
import { MAT_DIALOG_DATA } from '@angular/material/dialog';

@Component({
  selector: 'app-result',
  templateUrl: './result.component.html',
  styleUrls: ['./result.component.scss'],
})
export class ResultComponent implements OnInit {
  constructor(@Inject(MAT_DIALOG_DATA) public data: any) {}
  porcentaje = 0;
  ngOnInit(): void {
    this.porcentaje = parseFloat(this.data.toSend.prediction.slice(1, -1));
    this.porcentaje = this.porcentaje * 100;
    this.porcentaje = parseFloat(this.porcentaje.toFixed(2));
  }
}
