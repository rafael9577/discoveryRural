import { createStackNavigator } from '@react-navigation/stack'

import login from '../screens/login'
import signUp from '../screens/signup'
import tabsRouts from './TabsNavigator'
import infoAnimal from "../screens/userPage/infoAnimal"

const stack = createStackNavigator();

export default () => (
    <stack.Navigator initialRouteName='login' screenOptions={{ headerShown: false }}>
        <stack.Screen name='login' component={login} />
        <stack.Screen name='signup' component={signUp} />
        <stack.Screen name='tabsRouts' component={tabsRouts}/>
        <stack.Screen name='infoAnimal' component={infoAnimal} />
    </stack.Navigator>
)