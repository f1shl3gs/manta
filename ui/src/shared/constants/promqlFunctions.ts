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
    desc: 'abs() returns the input vector with all sample values converted to their absolute value',
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
    desc: 'absent(v instant-vector) returns an empty vector if the vector passed to it has any elements and a 1-element vector with the value 1 if the vector passed to it has no elements.',
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
  {
    name: 'ceil',
    args: [
      {
        name: 'v',
        type: 'instant-vector',
        desc: 'The input vector',
      },
    ],
    desc: 'ceil(v instant-vector) rounds the sample values of all elements in v up to the nearest integer',
    example: '',
  },
  {
    name: 'changes',
    args: [
      {
        name: 'v',
        type: 'instant-vector',
        desc: 'The input vector',
      },
    ],
    desc: 'For each input time series, changes(v range-vector) returns the number of times its value has changed within the provided time range as an instant vector.',
    example: 'changes(up{instance="127.0.0.1:9100"}[5m])',
  },
  {
    name: 'clamp_max',
    args: [
      {
        name: 'v',
        type: 'instant-vector',
        desc: 'The input vector',
      },
      {
        name: 'max',
        type: 'scalar',
        desc: 'max',
      },
    ],
    desc: 'clamp_max(v instant-vector, max scalar) clamps the sample values of all elements in v to have an upper limit of max.',
    example: '',
  },
  {
    name: 'clamp_min',
    args: [
      {
        name: 'v',
        type: 'instant-vector',
        desc: 'The input vector',
      },
      {
        name: 'min',
        type: 'scalar',
        desc: 'max',
      },
    ],
    desc: 'clamp_min(v instant-vector, min scalar) clamps the sample values of all elements in v to have a lower limit of min.',
    example: '',
  },
  {
    name: 'day_of_month',
    args: [],
    desc: 'day_of_month(v=vector(time()) instant-vector) returns the day of the month for each of the given times in UTC. Returned values are from 1 to 31',
    example: '',
  },
  {
    name: 'day_of_week',
    args: [],
    desc: 'day_of_week(v=vector(time()) instant-vector) returns the day of the week for each of the given times in UTC. Returned values are from 0 to 6, where 0 means Sunday etc',
    example: '',
  },
  {
    name: 'days_in_month',
    args: [],
    desc: 'days_in_month(v=vector(time()) instant-vector) returns number of days in the month for each of the given times in UTC. Returned values are from 28 to 31.',
    example: '',
  },
  {
    name: 'delta',
    args: [
      {
        name: 'v',
        type: 'range-vector',
      },
    ],
    desc: 'delta(v range-vector) calculates the difference between the first and last value of each time series element in a range vector v, returning an instant vector with the given deltas and equivalent labels. The delta is extrapolated to cover the full time range as specified in the range vector selector, so that it is possible to get a non-integer result even if the sample values are all integers.\n\ndelta should only be used with gauges.',
    example: 'delta(cpu_temp_celsius{host="zeus"}[2h])',
  },
  {
    name: 'deriv',
    args: [
      {
        name: 'v',
        type: 'range-vector',
      },
    ],
    desc:
      'deriv(v range-vector) calculates the per-second derivative of the time series in a range vector v, using simple linear regression.\n' +
      '\n' +
      'deriv should only be used with gauges.',
    example: '',
  },
  {
    name: 'exp',
    args: [
      {
        name: 'v',
        type: 'instant-vector',
      },
    ],
    desc:
      'exp(v instant-vector) calculates the exponential function for all elements in v. Special cases are:\n' +
      '\n' +
      'Exp(+Inf) = +Inf\n' +
      'Exp(NaN) = NaN\n',
    example: '',
  },
  {
    name: 'floor',
    args: [
      {
        name: 'v',
        type: 'instant-vector',
      },
    ],
    desc: 'floor(v instant-vector) rounds the sample values of all elements in v down to the nearest integer.',
    example: '',
  },
  {
    name: 'histogram_quantile',
    args: [
      {
        name: 'p',
        type: 'scaler',
        desc: '',
      },
      {
        name: 'b',
        type: 'instant-vector',
      },
    ],
    desc:
      'histogram_quantile(φ scalar, b instant-vector) calculates the φ-quantile (0 ≤ φ ≤ 1) from the buckets b of a histogram. (See histograms and summaries for a detailed explanation of φ-quantiles and the usage of the histogram metric type in general.) The samples in b are the counts of observations in each bucket. Each sample must have a label le where the label value denotes the inclusive upper bound of the bucket. (Samples without such a label are silently ignored.) The histogram metric type automatically provides time series with the _bucket suffix and the appropriate labels.\n' +
      '\n' +
      'Use the rate() function to specify the time window for the quantile calculation.',
    example:
      '# A histogram metric is called http_request_duration_seconds. To calculate the 90th percentile of request durations over the last 10m, use the following expression\n' +
      'histogram_quantile(0.9, rate(http_request_duration_seconds_bucket[10m]))',
    link: 'https://prometheus.io/docs/prometheus/latest/querying/functions/#histogram_quantile',
  },
  {
    name: 'holt_winters',
    args: [
      {
        name: 'v',
        type: 'range-vector',
        desc: 'The input vector',
      },
      {
        name: 'sf',
        type: 'scalar',
        desc: 'The lower sf the more importance of old data, sf must be between 0 and 1',
      },
      {
        name: 'tf',
        type: 'scalar',
        desc: 'The higher the more trends is considered, tf must be between 0 and 1',
      },
    ],
    desc:
      'holt_winters(v range-vector, sf scalar, tf scalar) produces a smoothed value for time series based on the range in v. The lower the smoothing factor sf, the more importance is given to old data. The higher the trend factor tf, the more trends in the data is considered. Both sf and tf must be between 0 and 1\n\n' +
      'holt_winters should only be used with gauges',
    example: '',
    link: 'https://prometheus.io/docs/prometheus/latest/querying/functions/#holt_winters',
  },
  {
    name: 'hour',
    args: [],
    desc: 'hour(v=vector(time()) instant-vector) returns the hour of the day for each of the given times in UTC. Returned values are from 0 to 23',
    example: '',
    link: 'https://prometheus.io/docs/prometheus/latest/querying/functions/#hour',
  },
  {
    name: 'idelta',
    args: [
      {
        name: 'v',
        type: 'range-vector',
        desc: 'The input vector',
      },
    ],
    desc:
      'idelta(v range-vector) calculates the difference between the last two samples in the range vector v, returning an instant vector with the given deltas and equivalent labels.\n' +
      '\n' +
      'idelta should only be used with gauges.',
    example: '',
    link: 'https://prometheus.io/docs/prometheus/latest/querying/functions/#idelta',
  },
  {
    name: 'increase',
    args: [
      {
        name: 'v',
        type: 'range-vector',
        desc: 'The input vector',
      },
    ],
    desc: 'increase(v range-vector) calculates the increase in the time series in the range vector. Breaks in monotonicity (such as counter resets due to target restarts) are automatically adjusted for. The increase is extrapolated to cover the full time range as specified in the range vector selector, so that it is possible to get a non-integer result even if a counter increases only by integer increments.',
    example: 'increase(http_requests_total{job="api-server"}[5m])',
    link: 'https://prometheus.io/docs/prometheus/latest/querying/functions/#increase',
  },
  {
    name: 'irate',
    args: [
      {
        name: 'v',
        type: 'range-vector',
        desc: 'The input vector',
      },
    ],
    desc: 'irate(v range-vector) calculates the per-second instant rate of increase of the time series in the range vector. This is based on the last two data points. Breaks in monotonicity (such as counter resets due to target restarts) are automatically adjusted for.',
    example: 'irate(http_requests_total{job="api-server"}[5m])',
    link: 'https://prometheus.io/docs/prometheus/latest/querying/functions/#irate',
  },
  {
    name: 'label_join',
    args: [
      {
        name: 'v',
        type: 'instant-vector',
        desc: 'The input vector',
      },
      {
        name: 'dst_label',
        type: 'string',
        desc: 'Dest label name',
      },
      {
        name: 'separator',
        type: 'string',
        desc: 'Separator of source labels',
      },
      {
        name: 'src_label_1',
        type: 'string',
        desc: 'Source label to join',
      },
      {
        name: 'src_label_2',
        type: 'string',
        desc: 'Source label to join',
      },
      {
        name: '...',
        type: 'string',
        desc: 'More source labels',
      },
    ],
    desc: 'For each timeseries in v, label_join(v instant-vector, dst_label string, separator string, src_label_1 string, src_label_2 string, ...) joins all the values of all the src_labels using separator and returns the timeseries with the label dst_label containing the joined value. There can be any number of src_labels in this function.',
    example:
      'label_join(up{job="api-server",src1="a",src2="b",src3="c"}, "foo", ",", "src1", "src2", "src3")',
    link: 'https://prometheus.io/docs/prometheus/latest/querying/functions/#label_join',
  },
  {
    name: 'label_replace',
    args: [
      {
        name: 'v',
        type: 'instant-vector',
        desc: 'The input vector',
      },
      {
        name: 'dst_label',
        type: 'string',
        desc: 'Dest label',
      },
      {
        name: 'replacement',
        type: 'string',
        desc: '',
      },
      {
        name: 'src_label',
        type: 'string',
        desc: 'Source label',
      },
      {
        name: 'regex',
        type: 'string',
        desc: 'Regular expression against the src_label',
      },
    ],
    desc: "For each timeseries in v, label_replace(v instant-vector, dst_label string, replacement string, src_label string, regex string) matches the regular expression regex against the label src_label. If it matches, then the timeseries is returned with the label dst_label replaced by the expansion of replacement. $1 is replaced with the first matching subgroup, $2 with the second etc. If the regular expression doesn't match then the timeseries is returned unchanged.",
    example:
      'label_replace(up{job="api-server",service="a:c"}, "foo", "$1", "service", "(.*):.*")',
    link: 'https://prometheus.io/docs/prometheus/latest/querying/functions/#label_replace',
  },
  {
    name: 'ln',
    args: [
      {
        name: 'v',
        type: 'instant-vector',
        desc: 'The input vector',
      },
    ],
    desc:
      'ln(v instant-vector) calculates the natural logarithm for all elements in v. Special cases are:\n' +
      '\n' +
      'ln(+Inf) = +Inf\n' +
      'ln(0) = -Inf\n' +
      'ln(x < 0) = NaN\n' +
      'ln(NaN) = NaN',
    example: '',
    link: 'https://prometheus.io/docs/prometheus/latest/querying/functions/#ln',
  },
  {
    name: 'log2',
    args: [
      {
        name: 'v',
        type: 'instant-vector',
        desc: 'The input vector',
      },
    ],
    desc: 'log2(v instant-vector) calculates the binary logarithm for all elements in v. The special cases are equivalent to those in ln',
    example: '',
    link: 'https://prometheus.io/docs/prometheus/latest/querying/functions/#log2',
  },
  {
    name: 'log10',
    args: [
      {
        name: 'v',
        type: 'instant-vector',
        desc: 'The input vector',
      },
    ],
    desc: 'log10(v instant-vector) calculates the decimal logarithm for all elements in v. The special cases are equivalent to those in ln.',
    example: '',
    link: 'https://prometheus.io/docs/prometheus/latest/querying/functions/#log10',
  },
  {
    name: 'minute',
    args: [],
    desc: 'minute(v=vector(time()) instant-vector) returns the minute of the hour for each of the given times in UTC. Returned values are from 0 to 59.',
    example: '',
    link: 'https://prometheus.io/docs/prometheus/latest/querying/functions/#minute',
  },
  {
    name: 'month',
    args: [],
    desc: 'month(v=vector(time()) instant-vector) returns the month of the year for each of the given times in UTC. Returned values are from 1 to 12, where 1 means January etc.',
    example: '',
    link: 'https://prometheus.io/docs/prometheus/latest/querying/functions/#month',
  },
  {
    name: 'predict_linear',
    args: [
      {
        name: 'v',
        type: 'range-vector',
        desc: 'The input vector',
      },
      {
        name: 't',
        type: 'scalar',
        desc: 'Seconds from now',
      },
    ],
    desc:
      'predict_linear(v range-vector, t scalar) predicts the value of time series t seconds from now, based on the range vector v, using simple linear regression.\n\n' +
      'predict_linear should only be used with gauges.',
    example: '',
    link: 'https://prometheus.io/docs/prometheus/latest/querying/functions/#predict_linear',
  },
  {
    name: 'rate',
    args: [
      {
        name: 'v',
        type: 'range-vector',
        desc: 'The input vector',
      },
    ],
    desc: "rate(v range-vector) calculates the per-second average rate of increase of the time series in the range vector. Breaks in monotonicity (such as counter resets due to target restarts) are automatically adjusted for. Also, the calculation extrapolates to the ends of the time range, allowing for missed scrapes or imperfect alignment of scrape cycles with the range's time period.",
    example: 'rate(http_requests_total{job="api-server"}[5m])',
    link: 'https://prometheus.io/docs/prometheus/latest/querying/functions/#rate',
  },
  {
    name: 'resets',
    args: [
      {
        name: 'v',
        type: 'range-vector',
        desc: 'The input vector',
      },
    ],
    desc:
      'For each input time series, resets(v range-vector) returns the number of counter resets within the provided time range as an instant vector. Any decrease in the value between two consecutive samples is interpreted as a counter reset.\n\n' +
      'resets should only be used with counters.',
    example: '',
    link: 'https://prometheus.io/docs/prometheus/latest/querying/functions/#resets',
  },
  {
    name: 'round',
    args: [
      {
        name: 'v',
        type: 'range-vector',
        desc: 'The input vector',
      },
      {
        name: 'to_nearest',
        type: 'scalar',
        desc: 'The nearest integer to round',
      },
    ],
    desc: 'round(v instant-vector, to_nearest=1 scalar) rounds the sample values of all elements in v to the nearest integer. Ties are resolved by rounding up. The optional to_nearest argument allows specifying the nearest multiple to which the sample values should be rounded. This multiple may also be a fraction.',
    example: '',
    link: 'https://prometheus.io/docs/prometheus/latest/querying/functions/#round',
  },
  {
    name: 'scalar',
    args: [
      {
        name: 'v',
        type: 'instant-vector',
        desc: 'The input vector',
      },
    ],
    desc: 'Given a single-element input vector, scalar(v instant-vector) returns the sample value of that single element as a scalar. If the input vector does not have exactly one element, scalar will return NaN.',
    example: '',
    link: 'https://prometheus.io/docs/prometheus/latest/querying/functions/#scalar',
  },
  {
    name: 'sort',
    args: [
      {
        name: 'v',
        type: 'instant-vector',
        desc: 'The input vector',
      },
    ],
    desc: 'sort(v instant-vector) returns vector elements sorted by their sample values, in ascending order.',
    example: '',
    link: 'https://prometheus.io/docs/prometheus/latest/querying/functions/#sort',
  },
  {
    name: 'sort_desc',
    args: [
      {
        name: 'v',
        type: 'instant-vector',
        desc: 'The input vector',
      },
    ],
    desc: 'sort(v instant-vector) returns vector elements sorted by their sample values, in descending order.',
    example: '',
    link: 'https://prometheus.io/docs/prometheus/latest/querying/functions/#sort_desc',
  },
  {
    name: 'sqrt',
    args: [
      {
        name: 'v',
        type: 'instant-vector',
        desc: 'The input vector',
      },
    ],
    desc: 'sqrt(v instant-vector) calculates the square root of all elements in v',
    example: '',
    link: 'https://prometheus.io/docs/prometheus/latest/querying/functions/#sqrt',
  },
  {
    name: 'time',
    args: [],
    desc: 'time() returns the number of seconds since January 1, 1970 UTC. Note that this does not actually return the current time, but the time at which the expression is to be evaluated.',
    example: '',
    link: 'https://prometheus.io/docs/prometheus/latest/querying/functions/#time',
  },
  {
    name: 'timestamp',
    args: [],
    desc: 'timestamp(v instant-vector) returns the timestamp of each of the samples of the given vector as the number of seconds since January 1, 1970 UTC.',
    example: '',
    link: 'https://prometheus.io/docs/prometheus/latest/querying/functions/#timestamp',
  },
  {
    name: 'vector',
    args: [
      {
        name: 's',
        type: 'scalar',
        desc: '',
      },
    ],
    desc: 'vector(s scalar) returns the scalar s as a vector with no labels.',
    example: '',
    link: 'https://prometheus.io/docs/prometheus/latest/querying/functions/#vector',
  },
  {
    name: 'year',
    args: [],
    desc: 'year(v=vector(time()) instant-vector) returns the year for each of the given times in UTC.',
    example: '',
    link: 'https://prometheus.io/docs/prometheus/latest/querying/functions/#year',
  },

  // aggregate operations
  {
    name: 'avg',
    args: [
      {
        name: 'v',
        type: 'range-vector',
        desc: 'The input vector',
      },
    ],
    desc: 'Calculate the average over dimensions',
    example: '',
    link: 'https://prometheus.io/docs/prometheus/latest/querying/operators/#aggregation-operators',
  },
  {
    name: 'bottomk',
    args: [
      {
        name: 'k',
        type: 'scalar',
        desc: 'The rank number',
      },
      {
        name: 'v',
        type: 'range-vector',
        desc: 'The input vector',
      },
    ],
    desc: 'Smallest k elements by sample value',
    example: '',
    link: 'https://prometheus.io/docs/prometheus/latest/querying/operators/#aggregation-operators',
  },
  {
    name: 'count',
    args: [
      {
        name: 'v',
        type: 'instant-vector | range-vector',
        desc: 'The input vector',
      },
    ],
    desc: 'Count number of elements in the vector',
    example: '',
    link: 'https://prometheus.io/docs/prometheus/latest/querying/operators/#aggregation-operators',
  },
  {
    name: 'count_values',
    args: [
      {
        name: 'label_name',
        type: 'scalar',
        desc: 'The label name',
      },
      {
        name: 'v',
        type: 'instant-vector',
        desc: 'The input vector',
      },
    ],
    desc: 'Count the number of series per distinct sample value',
    example: 'count_values("instance", up)',
    link: '',
  },
  /*{
    name: 'group',
    link:
      'https://prometheus.io/docs/prometheus/latest/querying/operators/#aggregation-operators',
  },*/
  {
    name: 'max',
    args: [
      {
        name: 'v',
        type: 'instant-vector | range-vector',
        desc: 'The input vector',
      },
    ],
    desc: 'Select maximum over dimensions',
    example: '',
    link: 'https://prometheus.io/docs/prometheus/latest/querying/operators/#aggregation-operators',
  },
  {
    name: 'min',
    args: [
      {
        name: 'v',
        type: 'instant-vector | range-vector',
        desc: 'The input vector',
      },
    ],
    desc: 'Select minimum over dimensions',
    example: '',
    link: 'https://prometheus.io/docs/prometheus/latest/querying/operators/#aggregation-operators',
  },
  {
    name: 'quantile',
    args: [
      {
        name: 'p',
        type: 'scalar',
        desc: 'The percentile',
      },
      {
        name: '...',
        type: 'instant-vector | range-vector',
        desc: 'The input vector',
      },
    ],
    desc: 'Calculate φ-quantile (0 ≤ φ ≤ 1) over dimensions',
    example: 'quantile(0.95, ...)',
    link: 'https://prometheus.io/docs/prometheus/latest/querying/operators/#aggregation-operators',
  },
  {
    name: 'stddev',
    args: [
      {
        name: 'v',
        type: 'instant-vector',
        desc: 'The input vector',
      },
    ],
    desc: 'Calculate population standard deviation over dimensions',
    example: '',
    link: 'https://prometheus.io/docs/prometheus/latest/querying/operators/#aggregation-operators',
  },
  {
    name: 'stdvar',
    args: [
      {
        name: 'v',
        type: 'instant-vector',
        desc: 'The input vector',
      },
    ],
    desc: 'Calculate population standard variance over dimensions',
    example: '',
    link: 'https://prometheus.io/docs/prometheus/latest/querying/operators/#aggregation-operators',
  },
  {
    name: 'sum',
    args: [
      {
        name: 'v',
        type: 'instant-vector',
        desc: 'The input vector',
      },
    ],
    desc: 'Calculate sum over dimensions',
    example: '',
    link: 'https://prometheus.io/docs/prometheus/latest/querying/operators/#aggregation-operators',
  },
  {
    name: 'topk',
    args: [
      {
        name: 'k',
        type: 'scalar',
        desc: 'The rank number',
      },
      {
        name: 'v',
        type: 'instant-vector',
        desc: 'The input vector',
      },
    ],
    desc: 'Largest k elements by sample value',
    example: '',
    link: 'https://prometheus.io/docs/prometheus/latest/querying/operators/#aggregation-operators',
  },
]
