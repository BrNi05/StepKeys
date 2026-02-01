// Strcture of a pedal action
// The same as defined on the backend
export interface PedalAction {
  mode: 'sequence' | 'combo';
  behaviour: 'oneshot' | 'toggle' | 'hold';
  keys: string[];
}
