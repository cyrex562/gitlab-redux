import { execute } from './executor';
import { parse } from './parser';
import { present } from './presenter';
import { transform } from './transformer/data';

export const executeAndPresentQuery = async (glqlQuery) => {
  const { query, config } = await parse(glqlQuery);
  const data = await execute(query);
  const transformed = transform(data, config);
  return present(transformed, config);
};

export const presentPreview = async (glqlQuery) => {
  const { config } = await parse(glqlQuery);
  const data = { project: { issues: { nodes: [] } } };
  const transformed = transform(data, config);
  return present(transformed, config, { isPreview: true });
};
