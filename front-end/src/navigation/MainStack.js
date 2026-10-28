import { createNativeStackNavigator } from '@react-navigation/native-stack'

import login from '../screens/login'
import signUp from '../screens/signup'
import tabsRouts from './TabsNavigator'
import infoAnimal from "../screens/userPage/infoAnimal"
import infoImportantes from "../screens/userPage/infoImportantes"
import animalDetalhe from "../screens/userPage/detalheAnimal"

const stack = createNativeStackNavigator();

export default () => (
    <stack.Navigator initialRouteName='login' screenOptions={{ headerShown: false }}>
        <stack.Screen name='login' component={login} />
        <stack.Screen name='signup' component={signUp} />
        <stack.Screen name='tabsRouts' component={tabsRouts}/>
        <stack.Screen name='infoAnimal' component={infoAnimal} />
        <stack.Screen name='infoImportantes' component={infoImportantes} />
        <stack.Screen name='animalDetalhe' component={animalDetalhe} />
    </stack.Navigator>
)