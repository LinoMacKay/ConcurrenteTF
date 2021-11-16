import { Injectable } from '@angular/core';
import { Symptom } from 'src/app/models/symptom';

@Injectable({
  providedIn: 'root',
})
export class SymptomsService {
  public symptoms = [
    new Symptom(
      'Tos',
      'https://www.laboratoriochile.cl/wp-content//uploads/2020/04/Hombre-tosiendo-1200-.min_.jpg',
      '',
      'tos'
    ),
    new Symptom(
      'Cefalea',
      'https://d2lcsjo4hzzyvz.cloudfront.net/blog/wp-content/uploads/2021/05/10180044/%C2%BFQue%CC%81-es-la-cefalea.jpg',
      'Dolor de Cabeza',
      'dolorcabeza'
    ),
    new Symptom(
      'Congestion Nasal',
      'https://geomedica.com.ar/wp-content/uploads/Congestio%CC%81n-nasal.2-1170x767.jpg',
      '',
      'congnasal'
    ),
    new Symptom(
      'Dificultad Respiratoria',
      'https://img.saludsavia.com/wp-content/uploads/2019/04/Dificultad-para-Respirar.jpg',
      '',
      'difresp'
    ),
    new Symptom(
      'Dolor de garganta',
      'https://f.elconfidencial.com/original/2aa/184/ea8/2aa184ea838571e5b307d2e9007a0984.jpg',
      '',
      'dolorgarg'
    ),
    new Symptom(
      'Fiebre',
      'https://cloudfront-us-east-1.images.arcpublishing.com/eluniverso/RPOFIBPC7BBSPKRJ57TUIHV37Y.png',
      'Temperatura corporal alrededor de los 37 °C',
      'fiebre'
    ),
    new Symptom(
      'Diarrea',
      'https://medlineplus.gov/images/Diarrhea.jpg',
      '',
      'diarrea'
    ),
    new Symptom(
      'Nauseas',
      'https://www.clikisalud.net/wp-content/uploads/2019/11/factores-nauseas-despues-comer.jpg',
      '',
      'nauseas'
    ),
    new Symptom(
      'Anosmia / Hiposmia',
      'https://www.news-medical.net/image-handler/ts/20210419064338/ri/673/picture/2021/4/shutterstock_1927048019.jpg',
      'Pérdida del sentido del gusto y olfato',
      'nooftgus'
    ),
    new Symptom(
      'Dolor Abdominal',
      'https://estaticos.muyinteresante.es/media/cache/1140x_thumb/uploads/images/report/5f367dca5bafe8cbe148778e/dolor-abdominal1_0.jpg',
      '',
      'dolorab'
    ),
    new Symptom(
      'Dolor de Articulaciones',
      'https://biosporty.com/wp-content/uploads/2019/09/dolor-articulaciones-640x360-1024x585.jpg',
      '',
      'dolorart'
    ),
    new Symptom(
      'Dolor Muscular',
      'https://static3.abc.es/media/bienestar/2020/09/18/dolor-muscular-kpfF--620x349@abc.jpg',
      '',
      'dolormusc'
    ),
    new Symptom(
      'Dolor de Pecho',
      'https://www.webconsultas.com/sites/default/files/styles/wc_adaptive_image__small/public/migrated/sintomas-angina-pecho.jpg',
      '',
      'dolorpech'
    ),
    new Symptom(
      'Otros Sintomas',
      'https://estaticos.muyinteresante.es/media/cache/1140x_thumb/uploads/images/article/5ea9253f5cafe82a84745a1f/nuevos-sintomas_0.jpg',
      '',
      'otros'
    ),
  ];

  constructor() {}

  selectSymptom(symptom: Symptom) {
    var index = this.symptoms.indexOf(symptom);
    this.symptoms[index].isSelected = !this.symptoms[index].isSelected;
  }

  GetSymtoms() {
    return this.symptoms;
  }

  getOnlySelected() {
    var symptonsToSend: { sintoma: string; isSelected: number }[] = [];

    this.symptoms.forEach((e) => {
      var data = {
        sintoma: e.shortName,
        isSelected: e.isSelected ? 1 : 0,
      };
      symptonsToSend.push(data);
    });

    return symptonsToSend;
  }
}
