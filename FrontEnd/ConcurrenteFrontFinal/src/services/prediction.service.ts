import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { environment } from 'src/environments/environment';

@Injectable({
  providedIn: 'root',
})
export class PredictionService {
  url: string = `${environment.HOST_URL}/predict`;

  constructor(private http: HttpClient) {}

  createPrediction(data: any) {
    return this.http.post<any>(`${this.url}`, data);
  }
}
