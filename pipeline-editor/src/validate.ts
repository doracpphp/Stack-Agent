import type { Pipeline, ValidationError } from './types';
import { getLocalizedTaskDef } from './taskDefs';
import type { Locale, T } from './i18n';

const ID_RE = /^[a-zA-Z0-9_]+$/;

export function validate(pipeline: Pipeline, t: T, locale: Locale): ValidationError[] {
  const errors: ValidationError[] = [];

  if (!pipeline.name.trim()) {
    errors.push({ message: t.valPipelineName });
  }

  const needsLLM = pipeline.tasks.some(
    (task) => task.type === 'llm_prompt' || task.type === 'fill_template',
  );
  if (needsLLM) {
    if (!pipeline.llm.provider) errors.push({ message: t.valLLMProvider });
    if (!pipeline.llm.model.trim()) errors.push({ message: t.valLLMModel });
  }

  const seenIds = new Set<string>();

  for (const task of pipeline.tasks) {
    if (!task.id.trim()) {
      errors.push({ taskKey: task._key, field: 'id', message: t.valTaskIdRequired });
    } else if (!ID_RE.test(task.id)) {
      errors.push({ taskKey: task._key, field: 'id', message: t.valTaskIdFormat });
    } else if (seenIds.has(task.id)) {
      errors.push({ taskKey: task._key, field: 'id', message: t.valTaskIdDuplicate(task.id) });
    }
    seenIds.add(task.id);

    const def = getLocalizedTaskDef(task.type, locale);
    if (!def) continue;

    for (const param of def.params) {
      if (!param.required) continue;
      const val = task.params[param.key];
      const empty =
        val === undefined ||
        val === null ||
        val === '' ||
        (Array.isArray(val) && val.length === 0);
      if (empty) {
        errors.push({ taskKey: task._key, field: param.key, message: t.valFieldRequired(param.label) });
      }
    }
  }

  return errors;
}
