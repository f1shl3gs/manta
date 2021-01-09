import {isEqual} from 'lodash';

describe('Grid', () => {
  it('isEqual', () => {
    const prev = {
      h: 4,
      i: '06dde20d9063d000',
      moved: false,
      static: false,
      w: 12,
      x: 0,
      y: 0
    };

    const next = {
      h: 4,
      i: '06dde20d9063d000',
      isBounded: undefined,
      isDraggable: undefined,
      isResizable: undefined,
      maxH: undefined,
      maxW: undefined,
      minH: undefined,
      minW: undefined,
      resizeHandles: undefined,
      moved: false,
      static: false,
      w: 12,
      x: 0,
      y: 0
    };

    console.log(isEqual(prev, next))
  });
});