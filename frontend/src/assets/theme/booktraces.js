import { definePreset } from '@primevue/themes'
import Aura from '@primevue/themes/aura'
import './styleoverrides.scss'

const BookTraces = definePreset(Aura, {
   root: {
      borderRadius: {
         none: '0',
         xs: '2px',
         sm: '3px',
         md: '4px',
         lg: '4px',
         xl: '8px'
      },
   },
   semantic: {
      primary: {
         50: '#24890d',
         100: '#24890d',
         200: '#24890d',
         300: '#24890d',
         400: '#24890d',
         500: '#24890d',
         600: '#24890d',
         700: '#24890d',
         800: '#24890d',
         900: '#24890d',
         950: '#24890d'
      },
      focusRing: {
         width: '1px',
         style: 'dashed',
         offset: '2px'
      },
      formField: {
         paddingX: '0.75rem',
         paddingY: '0.5rem',
         borderRadius: '4px',
         focusRing: {
             width: '1px',
             style: 'dashed',
             color: '{primary.color}',
             offset: '2px',
             shadow: 'none'
         },
         transitionDuration: '{transition.duration}'
      },
      disabledOpacity: '0.3',
      colorScheme: {
         light: {
            primary: {
               color: '{primary.500}',
               contrastColor: '#ffffff',
               hoverColor: '{primary.100}',
               activeColor: '{primary.500}'
            },
            highlight: {
               background: '#ffffff',
               focusBackground: '#ffffff',
               color: '#333',
               focusColor: '#ffffff'
            }
         },
      }
   },
   components: {
      button: {
         root: {
            paddingY: '.5em',
            paddingX: '1em',
         },
         colorScheme: {
            light: {
               secondary: {
                  background: '#f1f1f1',
                  hoverBackground: '#dadada',
                  hoverBorderColor: '808080',
                  borderColor: '#dadada',
                  color: '#333',
               },
            }
         }
      },
      datatable: {
         paginatorTop: {
            borderColor: '#f1f1f1',
         },
         headerCell: {
            borderColor: '#f1f1f1',
         },
         bodyCell: {
            borderColor: '#f1f1f1',
        },
         colorScheme: {
            light: {
               root: {
                  borderColor: 'transparent'
              },
            }
         }
      },
      dialog: {
         colorScheme: {
            light: {
               root: {
                  background: '#ffffff',
                  borderColor: '808080',
                  padding: '15px',
                  borderRadius: '4px',
               },
               header: {
                  padding: '10px',
               },
               title: {
                  fontWeight: '600',
                  fontSize: '1em',
               }
            }
         }
      },
      panel: {
         header: {
            background: '#f8f9fa',
            borderColor:  '#dadada',
            borderRadius: '4px 4px 0 0',
            padding: '1rem'
         },
         title: {
            fontWeight: '600',
         },
      },
      select: {
         root: {
            paddingY: '.5em',
            paddingX: '.5em',
            disabledBackground: '#fafafa',
            disabledColor: '#cacaca',
         },
         option: {
            selectedFocusBackground: '#bfe7f7',
            selectedFocusColor: '#333',
            selectedBackground: '#bfe7f7',
            selectedColor: '#333'
         }
      },
   },
});

export default BookTraces;