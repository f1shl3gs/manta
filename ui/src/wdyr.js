import React from 'react';

const enable = false;

if (process.env.NODE_ENV === 'development' && enable) {
  const whyDidYouRender = require('@welldone-software/why-did-you-render');
  whyDidYouRender(React, {
    trackAllPureComponents: true,
    trackHooks: true,
    include: [/^Button$/, /^Page/, /Header/, /Nav/, /Otcl.*/, /Org/],
  });
}
