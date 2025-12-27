import 'vuetify/styles'
import { createVuetify } from 'vuetify'
import { aliases, mdi } from 'vuetify/iconsets/mdi'
import colors from 'vuetify/util/colors'
import { VDataTable, VDatePicker, VTimePicker } from 'vuetify/components'

export default createVuetify({
  components: {
    VDataTable,
    VDatePicker,
    VTimePicker
  },
  theme: {
    defaultTheme: 'dark',
    themes: {
      light: {
        colors: {
          primary: colors.blue.lighten2,
          secondary: colors.grey.darken1,
          accent: colors.shades.black,
          error: colors.red.accent3
        }
      },
      dark: {
        colors: {
          primary: colors.blue.darken3
        }
      }
    }
  },
  icons: {
    defaultSet: 'mdi',
    aliases,
    sets: {
      mdi
    }
  }
})
