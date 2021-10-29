import { Component, OnInit } from '@angular/core';
import {
  Router,
  NavigationStart,
  Event as NavigationEvent,
  ActivationStart,
} from '@angular/router';
@Component({
  selector: 'app-appbar',
  templateUrl: './appbar.component.html',
  styleUrls: ['./appbar.component.scss'],
})
export class AppbarComponent implements OnInit {
  constructor(private router: Router) {}
  visibility = true;

  ngOnInit(): void {
    this.router.events.subscribe((event: NavigationEvent) => {
      if (event instanceof ActivationStart) {
        if (event.snapshot.data.appbar != null) {
          this.visibility = event.snapshot.data.appbar;
        }
      }
    });
  }
}
