export class Symptom {
  name: string = '';
  image: string = '';
  definition: string = '';
  isSelected = false;

  constructor(name: string, image: string, definition: string) {
    this.name = name;
    this.image = image;
    this.definition = definition;
    // TBD
  }
}
