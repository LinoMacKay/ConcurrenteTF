import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { HomeComponent } from './core/home/home.component';
import { AppbarComponent } from './shared/appbar/appbar.component';
import { PredictionComponent } from './core/prediction/prediction.component';
import { ReactiveFormsModule } from '@angular/forms';
import { SimtomasCardComponent } from './core/simtomas-card/simtomas-card.component';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { AngularMaterialModule } from 'src/material/material.module';
import { ResultComponent } from './core/result/result.component';
import { NotificationComponent } from './shared/notification/notification.component';

@NgModule({
  declarations: [
    AppComponent,
    HomeComponent,
    AppbarComponent,
    PredictionComponent,
    SimtomasCardComponent,
    ResultComponent,
    NotificationComponent,
  ],
  imports: [
    BrowserModule,
    AppRoutingModule,
    ReactiveFormsModule,
    BrowserAnimationsModule,
    AngularMaterialModule,
  ],
  providers: [],
  bootstrap: [AppComponent],
})
export class AppModule {}
