import { Feather } from '@expo/vector-icons';
import { createBottomTabNavigator} from '@react-navigation/bottom-tabs'

const tabs = createBottomTabNavigator();

import home from "../screens/userPage/home"
import busca from "../screens/userPage/informacao"

export default() => (
    <tabs.Navigator 
        initialRouteName='home' 
        screenOptions={({ route }) => ({
            headerShown: false,
            tabBarIcon: ({ color, size }) => {
                let iconName;
                if (route.name === 'home') {
                    iconName = 'home';
                } else if (route.name === 'profile') {
                    iconName = 'user';
                } else if (route.name === 'settings') {
                    iconName = 'settings';
                } else if (route.name === 'busca') {
                    iconName = 'search';
                }
                return <Feather name={iconName} size={size} color={color} />;
            },
            tabBarActiveTintColor: '#2E7D32',
            tabBarInactiveTintColor: 'gray' })}>

        <tabs.Screen name='home' component={home}/>
        <tabs.Screen name='busca' component={busca}/>

    </tabs.Navigator>
)
