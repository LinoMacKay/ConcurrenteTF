export class Symptom {
  name: string = '';
  image: string = '';
  definition: string = '';
  isSelected = false;
  shortName = '';

  constructor(
    name: string,
    image: string,
    definition: string,
    shortName: string
  ) {
    this.name = name;
    this.image = image;
    this.definition = definition;
    this.shortName = shortName;
    // TBD
  }
}
