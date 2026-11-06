import { createBottomTabNavigator } from '@react-navigation/bottom-tabs';
import { useNavigation } from "@react-navigation/native";
import { Feather } from '@expo/vector-icons';
import sair from '../screens/userPage/sair';

import home from "../screens/userPage/home";
import busca from "../screens/userPage/informacao";
import emergencia from '../screens/userPage/emergencia';
import mapa from '../screens/userPage/mapa';

const tabs = createBottomTabNavigator();

const TelaVazia = () => null;

export default () => {
    const navigation = useNavigation();

    return (
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
                    } else if (route.name === 'sair') {
                        iconName = 'log-out';
                    } else if (route.name === 'emergencia') {
                        iconName = 'help-circle';
                    }else if (route.name === 'mapa') {
                        iconName = 'map';
                    } 
                    return <Feather name={iconName} size={size} color={color} />;
                },
                tabBarActiveTintColor: '#2E7D32',
                tabBarInactiveTintColor: 'gray' 
            })}
        >
            <tabs.Screen name='home' component={home}/>
            <tabs.Screen name='busca' component={busca}/>
            <tabs.Screen name='mapa' component={mapa}/>
            <tabs.Screen name='emergencia' component={emergencia} />
            
            <tabs.Screen 
                name='sair' 
                component={TelaVazia} 
                listeners={{
                    tabPress: async (e) => {
                        e.preventDefault(); 
                        try {
                            await sair();     
                        } catch (error) {
                            console.log("Erro ao sair:", error);
                        } finally {
                            navigation.navigate("login");
                        }
                    },
                }}
            />
        </tabs.Navigator>
    );
};
