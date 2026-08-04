export interface Lab {
  id: number;
  name: string;
  location: string;
  capacity: number;
}

export type NewLab = Omit<Lab, 'id'>;
