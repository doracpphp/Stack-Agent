export type LLMProvider = 'ollama' | 'openai' | 'gemini';

export interface LLMConfig {
  provider: LLMProvider;
  model: string;
  base_url: string;
  api_key: string;
}

export interface Task {
  /** Unique key for React rendering (not exported to JSON) */
  _key: string;
  id: string;
  type: string;
  params: Record<string, unknown>;
}

export interface Pipeline {
  name: string;
  version: string;
  llm: LLMConfig;
  tasks: Task[];
}

export interface ValidationError {
  taskKey?: string;
  field?: string;
  message: string;
}
