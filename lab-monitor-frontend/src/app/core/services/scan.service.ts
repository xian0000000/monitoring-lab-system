import { HttpClient } from '@angular/common/http';
import { Injectable, OnDestroy, inject, signal } from '@angular/core';
import { Observable, Subject } from 'rxjs';
import { ScanResult } from '../models/device.model';

@Injectable({ providedIn: 'root' })
export class ScanService implements OnDestroy {
  private http = inject(HttpClient);
  private eventSource?: EventSource;
  private resultsSubject = new Subject<ScanResult>();

  readonly connected = signal(false);
  readonly results$ = this.resultsSubject.asObservable();

  connect(): void {
    if (this.eventSource) {
      return;
    }

    this.eventSource = new EventSource('/api/scan/stream');

    this.eventSource.onopen = () => this.connected.set(true);
    this.eventSource.onerror = () => this.connected.set(false);
    this.eventSource.onmessage = (event: MessageEvent<string>) => {
      const result: ScanResult = JSON.parse(event.data);
      this.resultsSubject.next(result);
    };
  }

  disconnect(): void {
    this.eventSource?.close();
    this.eventSource = undefined;
    this.connected.set(false);
  }

  triggerScan(): Observable<{ message: string }> {
    return this.http.post<{ message: string }>('/api/scan', {});
  }

  ngOnDestroy(): void {
    this.disconnect();
  }
}
