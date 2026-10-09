import { Feather } from '@expo/vector-icons';
import { createBottomTabNavigator} from '@react-navigation/bottom-tabs'

const tabs = createBottomTabNavigator();

import home from "../screens/home"

export default() => {
    <tabs.Navigator initialRouteName='home' screenOptions={false} >
        <tabs.Screen name='home' component={home} options={{ tabBarIcon: ({color, size}) => {<Feather name='' color={color} size={size}/>} }} />
    </tabs.Navigator>
}
