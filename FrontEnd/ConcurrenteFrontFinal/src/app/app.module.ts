import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { HomeComponent } from './core/home/home.component';
import { AppbarComponent } from './shared/appbar/appbar.component';
import { PredictionComponent } from './core/prediction/prediction.component';

@NgModule({
  declarations: [
    AppComponent,
    HomeComponent,
    AppbarComponent,
    PredictionComponent
  ],
  imports: [
    BrowserModule,
    AppRoutingModule
  ],
  providers: [],
  bootstrap: [AppComponent]
})
export class AppModule { }
