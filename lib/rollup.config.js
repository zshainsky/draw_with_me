import resolve from '@rollup/plugin-node-resolve';
//import { terser } from 'rollup-plugin-terser';

export default [
  {
    input: './components/rooms-grid.js',
    output: [
      {
        dir: '../static/js',
        format: 'iife',
        sourcemap: true,
      },
    ],
    plugins: [
      resolve({
        browser: true,
      }),
      // for minification
      //terser(),
    ],
  },{
    input: './components/styles.js',
    output: [
      {
        dir: '../static/js',
        format: 'iife',
        sourcemap: true,
      },
    ],
    plugins: [
      resolve({
        browser: true,
      }),
      // for minification
      //terser(),
    ],
  },{
    input: './components/navigation-bar.js',
    output: [
      {
        dir: '../static/js',
        format: 'iife',
        sourcemap: true,
      },
    ],
    plugins: [
      resolve({
        browser: true,
      }),
      // for minification
      //terser(),
    ],
  },{
    input: './components/room-element.js',
    output: [
      {
        dir: '../static/js',
        format: 'iife',
        sourcemap: true,
      },
    ],
    plugins: [
      resolve({
        browser: true,
      }),
      // for minification
      //terser(),
    ],
  },{
    input: './components/footer-bar.js',
    output: [
      {
        dir: '../static/js',
        format: 'iife',
        sourcemap: true,
      },
    ],
    plugins: [
      resolve({
        browser: true,
      }),
      // for minification
      //terser(),
    ],
  },{
    input: './components/room-canvas.js',
    output: [
      {
        dir: '../static/js',
        format: 'iife',
        sourcemap: true,
      },
    ],
    plugins: [
      resolve({
        browser: true,
      }),
      // for minification
      //terser(),
    ],
  },{
    input: './components/tool-palette.js',
    output: [
      {
        dir: '../static/js',
        format: 'iife',
        sourcemap: true,
      },
    ],
    plugins: [
      resolve({
        browser: true,
      }),
      // for minification
      //terser(),
    ],
  },{
    input: './components/active-user-bar.js',
    output: [
      {
        dir: '../static/js',
        format: 'iife',
        sourcemap: true,
      },
    ],
    plugins: [
      resolve({
        browser: true,
      }),
      // for minification
      //terser(),
    ],
  },
];
