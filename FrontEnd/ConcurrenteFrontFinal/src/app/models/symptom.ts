export class Symptom {
  name: string = '';
  image: string = '';
  isSelected = false;

  constructor(name: string, image: string) {
    this.name = name;
    this.image = image;
    // TBD
  }
}
