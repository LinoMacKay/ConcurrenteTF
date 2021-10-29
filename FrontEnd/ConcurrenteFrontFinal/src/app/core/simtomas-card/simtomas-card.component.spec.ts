import { ComponentFixture, TestBed } from '@angular/core/testing';

import { SimtomasCardComponent } from './simtomas-card.component';

describe('SimtomasCardComponent', () => {
  let component: SimtomasCardComponent;
  let fixture: ComponentFixture<SimtomasCardComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ SimtomasCardComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(SimtomasCardComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
