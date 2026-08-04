import { HttpClient } from '@angular/common/http';
import { Injectable, inject } from '@angular/core';
import { Observable } from 'rxjs';
import { Lab, NewLab } from '../models/lab.model';

@Injectable({ providedIn: 'root' })
export class LabService {
  private http = inject(HttpClient);
  private baseUrl = '/api/labs';

  getLabs(): Observable<Lab[]> {
    return this.http.get<Lab[]>(this.baseUrl);
  }

  createLab(lab: NewLab): Observable<Lab> {
    return this.http.post<Lab>(this.baseUrl, lab);
  }
}
