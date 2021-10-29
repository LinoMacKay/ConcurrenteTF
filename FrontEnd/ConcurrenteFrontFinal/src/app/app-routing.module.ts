import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { HomeComponent } from './core/home/home.component';
import { PredictionComponent } from './core/prediction/prediction.component';

const routes: Routes = [
  { path: '', redirectTo: '/prediction', pathMatch: 'full' },
  { path: 'home', component: HomeComponent, data: { appbar: false } },
  { path: 'prediction', component: PredictionComponent },
];

@NgModule({
  imports: [RouterModule.forRoot(routes, { onSameUrlNavigation: 'reload' })],
  exports: [RouterModule],
})
export class AppRoutingModule {}
