import React from 'react';
import { SafeAreaView, ScrollView, View, Text, TouchableOpacity, Linking } from 'react-native';
import { useNavigation } from '@react-navigation/native';
import { styles } from './styled';

export default function InfoImportantes() {
  const navigation = useNavigation();

  const handleCall = (numero) => {
    Linking.openURL(`tel:${numero}`);
  };

  return (
    <SafeAreaView style={styles.container}>
      <ScrollView contentContainerStyle={styles.scrollContent} showsVerticalScrollIndicator={false}>
        
        <View style={styles.header}>
          <TouchableOpacity style={styles.backButton} onPress={() => navigation.goBack()}>
            <Text style={styles.title}>←</Text>
          </TouchableOpacity>
          <View style={styles.headerTextContainer}>
            <Text style={styles.title}>Guia de Emergência</Text>
            <Text style={styles.subtitle}>Suporte para Recife e Região Metropolitana</Text>
          </View>
        </View>

        <View style={styles.section}>
          <Text style={styles.sectionTitle}>Primeiros Socorros e Protocolos</Text>
          
          <View style={styles.card}>
            <Text style={styles.cardTitle}>Em caso de picada ou mordida</Text>
            <Text style={styles.cardText}><Text style={styles.highlightText}>O que fazer:</Text> Lave o local apenas com água e sabão neutro. Mantenha a vítima calma, deitada e hidratada. Retire anéis ou calçados próximos à área atingida.</Text>
            <Text style={styles.cardText}><Text style={styles.highlightText}>O que NÃO fazer:</Text> Não faça torniquete, não corte ou fure o local, não tente sugar o veneno e não aplique substâncias caseiras.</Text>
            <Text style={styles.cardText}><Text style={styles.highlightText}>Dica:</Text> Se for seguro, tire uma foto do animal à distância para apresentar à equipe médica (ajuda na escolha do soro).</Text>
          </View>
          
          <View style={styles.card}>
            <Text style={styles.cardTitle}>Ao avistar um animal silvestre</Text>
            <Text style={styles.cardText}>Mantenha distância mínima de segurança.</Text>
            <Text style={styles.cardText}>Não tente alimentar, acariciar ou capturar (Saguis, Iguanas, Serpentes, Capivaras).</Text>
            <Text style={styles.cardText}>Mantenha animais domésticos (cães/gatos) afastados dos animais silvestres.</Text>
          </View>
        </View>

        <View style={styles.section}>
          <Text style={styles.sectionTitle}>Contatos de Emergência (PE)</Text>
          
          <View style={styles.card}>
            <Text style={styles.cardTitle}>CEATOX-PE (Assistência Toxicológica)</Text>
            <Text style={styles.cardText}>Orientação médica 24h para acidentes com animais peçonhentos.</Text>
            <TouchableOpacity onPress={() => handleCall('08007226001')}>
              <Text style={styles.cardText}><Text style={styles.highlightText}>0800 722 6001</Text></Text>
            </TouchableOpacity>
          </View>
          
          <View style={styles.card}>
            <Text style={styles.cardTitle}>CIPOMA / BPMA (Polícia Ambiental)</Text>
            <Text style={styles.cardText}>Resgate e recolhimento de animais silvestres no Grande Recife.</Text>
            <TouchableOpacity onPress={() => handleCall('190')}>
              <Text style={styles.cardText}><Text style={styles.highlightText}>(81) 3181-1700 / 190</Text></Text>
            </TouchableOpacity>
          </View>
          
          <View style={styles.card}>
            <Text style={styles.cardTitle}>Corpo de Bombeiros / SAMU</Text>
            <TouchableOpacity onPress={() => handleCall('193')}>
              <Text style={styles.cardText}>Bombeiros (Resgate de risco): <Text style={styles.highlightText}>193</Text></Text>
            </TouchableOpacity>
            <TouchableOpacity onPress={() => handleCall('192')}>
              <Text style={styles.cardText}>SAMU (Urgência sistêmica grave): <Text style={styles.highlightText}>192</Text></Text>
            </TouchableOpacity>
          </View>
        </View>

        <View style={styles.footer}>
          {/* <Text style={styles.footerText}>A preservação da fauna nativa e dos ecossistemas de Pernambuco depende da nossa convivência consciente.</Text>*/}
        </View>

      </ScrollView>
    </SafeAreaView>
  );
}