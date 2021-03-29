export interface PromqlFunction {
  name: string
  args: {
    name: string
    type: string
    desc?: string
  }[]
  desc: string
  example: string
  link?: string
}

export const PROMQL_FUNCTIONS: PromqlFunction[] = [
  {
    name: 'abs',
    args: [
      {
        name: 'v',
        type: 'instant-vector',
        desc: 'The input vector.',
      },
    ],
    desc:
      'abs() returns the input vector with all sample values converted to their absolute value',
    example: 'abs(up{job="node_exporter"})',
  },
  {
    name: 'absent',
    args: [
      {
        name: 'v',
        type: 'instant-vector',
        desc: 'The input vector.',
      },
    ],
    desc:
      'absent(v instant-vector) returns an empty vector if the vector passed to it has any elements and a 1-element vector with the value 1 if the vector passed to it has no elements.',
    example:
      'absent(nonexistent{job="myjob"})\n' +
      '# => {job="myjob"}\n' +
      '\n' +
      'absent(nonexistent{job="myjob",instance=~".*"})\n' +
      '# => {job="myjob"}\n' +
      '\n' +
      'absent(sum(nonexistent{job="myjob"}))\n' +
      '# => {}',
  },
  {
    name: 'absent_over_time',
    args: [
      {
        name: 'v',
        type: 'instant-vector',
        desc: 'The input vector',
      },
    ],
    desc:
      'absent_over_time(v range-vector) returns an empty vector if the range vector passed to it has any elements and a 1-element vector with the value 1 if the range vector passed to it has no elements.\n' +
      '\n' +
      'This is useful for alerting on when no time series exist for a given metric name and label combination for a certain amount of time.',
    example:
      'absent_over_time(nonexistent{job="myjob"}[1h])\n' +
      '# => {job="myjob"}\n' +
      '\n' +
      'absent_over_time(nonexistent{job="myjob",instance=~".*"}[1h])\n' +
      '# => {job="myjob"}\n' +
      '\n' +
      'absent_over_time(sum(nonexistent{job="myjob"})[1h:])\n' +
      '# => {}',
  },
]
